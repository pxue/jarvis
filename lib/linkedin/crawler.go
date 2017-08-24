package linkedin

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goware/lg"
)

type WorkExperience struct {
	Company   string    `json:"company"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type WorkProfile struct {
	HashID     string            `json:"id"`
	XPs        []*WorkExperience `json:"experiences"`
	profileURL string
}

const (
	shortForm = "January 2006"
	yearForm  = "2006"
)

func (c *Client) CrawlProfiles() error {
	f, err := os.OpenFile(fmt.Sprintf("%s/%s_churn.json", DataDir, c.Company), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	parsed, err := c.LoadParsedSearch()
	if err != nil {
		return err
	}

	t := false
	for _, cn := range parsed {
		if cn.FullName == "SaadKabir" || t {
			t = true
		} else {
			continue
		}

		xurl, _ := url.Parse(cn.ProfileURL)
		xq, _ := url.ParseQuery(xurl.RawQuery)
		nq := url.Values{}
		nq.Set("id", xq.Get("id"))
		nq.Set("authType", xq.Get("authType"))
		nq.Set("authToken", xq.Get("authToken"))
		xurl.RawQuery = nq.Encode()

		res, err := c.Get(xurl.String())
		if err != nil {
			return err
		}

		hID := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(cn.FullName))))
		wp := &WorkProfile{HashID: hID}
		if err := wp.GetWorkExperience(res.Body); err != nil || len(wp.XPs) == 0 {
			lg.Printf("no xp found for %s (%s), skipping", cn.FullName, wp.profileURL)
		} else {
			lg.Printf("parsing %s, found %d xp", cn.FullName, len(wp.XPs))
			enc.Encode(wp)
		}

		res.Body.Close()
		c.Sleep()
	}

	return nil
}

func (w *WorkProfile) GetWorkExperience(r io.ReadCloser) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	defer r.Close()

	doc.Find("div[id^=experience][id$=view]").Each(func(_ int, s *goquery.Selection) {
		xp := &WorkExperience{
			Title:   s.Find("header h4").Text(),
			Company: s.Find("header h5").Text(),
		}

		s.Find(".experience-date-locale time").Each(func(i int, s *goquery.Selection) {
			t := shortForm
			if len(s.Text()) < 5 {
				t = yearForm
			}

			d, _ := time.Parse(t, s.Text())
			if i == 0 {
				xp.StartDate = d
			} else {
				xp.EndDate = d
			}
		})

		w.XPs = append(w.XPs, xp)
	})

	// try grabbing the profile url
	w.profileURL = doc.Find("dl.public-profile dd").Text()

	return nil
}
