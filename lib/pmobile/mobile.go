package mobile

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goware/lg"
	"github.com/pxue/jarvis/lib/form"
)

const (
	BaseUrl   = "https://selfserve.publicmobile.ca/"
	UserField = "ctl00$FullContent$ContentBottom$LoginControl$UserName"
	PassField = "ctl00$FullContent$ContentBottom$LoginControl$Password"
)

type Overview struct {
	DataAllowed float64
	DataUsed    float64
}

func ParseOverview(cookies []*http.Cookie) (*Overview, error) {
	req, _ := http.NewRequest("GET", BaseUrl+"Overview/", nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}

	res, _ := http.DefaultClient.Do(req)
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	ov := &Overview{}

	allowedStr := strings.TrimSpace(strings.TrimRight(doc.Find("span#VoiceAllowanceLiteral").Text(), "MB"))
	ov.DataAllowed, err = strconv.ParseFloat(allowedStr, 64)
	ov.DataUsed, _ = strconv.ParseFloat(doc.Find("span#VoiceUsedLiteral").Text(), 64)

	return ov, nil
}

func Login() error {
	// get cookie and stuff
	res, err := http.Get(BaseUrl)
	if err != nil {
		return err
	}

	// find the form
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return err
	}

	f := form.NewForm(BaseUrl, doc.Find("form"))
	f.SetInput(UserField, os.Getenv("PM_USERNAME"))
	f.SetInput(PassField, os.Getenv("PM_PASSWORD"))
	f.SetInput("__EVENTARGUMENT", "")
	f.SetInput("__EVENTTARGET", "")

	res, err = f.Submit()
	if err != nil {
		return err
	}

	for _, c := range res.Cookies() {
		log.Printf("%s: %s", c.Name, c.Value)
	}

	ov, _ := ParseOverview(res.Cookies())
	lg.Printf("%+v", ov)

	return nil
}
