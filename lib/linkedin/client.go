package linkedin

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/goware/lg"
)

type Client struct {
	Company   string
	CompanyID int64

	cl *http.Client

	// session cookies are injected every request
	// only `li_at` needs to be set
	cookies []*http.Cookie
	headers http.Header
	rand    *rand.Rand
}

const (
	ApiUrl  = "https://api.linkedin.com"
	DataDir = "./tmp"
)

func redirectPolFunc(r *http.Request, via []*http.Request) error {
	lg.Println(r.RequestURI)

	return nil
}

func NewClient(company string, companyID int64) (*Client, error) {
	//s := os.Getenv("LINKEDIN_SESSION")
	//if s == "" {
	//return nil, errors.New("linkedin env not set")
	//}
	s := "AQEDASCM_loBSKM5AAABWWJdffIAAAFZZBTx8k0Ahb2sIqfOrUVX1AJwIOdybaIuAC4Da6EmFUaUI5QhAgLt189z2PRDBtuv5V5tmzmqqw_yx4elgSchlcKnDK3IbpScArBKSYKW9tTWZjoxeA4N8y3I"
	return &Client{
		Company:   company,
		CompanyID: companyID,
		cl: &http.Client{
			CheckRedirect: redirectPolFunc,
		},
		cookies: []*http.Cookie{
			{Name: "li_at", Value: s},
		},
		headers: http.Header{
			"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36"},
		},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (c *Client) Get(getUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	return c.cl.Do(req)
}

func (c *Client) Sleep() {
	// sleep for random inteval
	n := 3 + c.rand.Intn(10)
	lg.Printf("sleeping for %d seconds", n)
	time.Sleep(time.Duration(n) * time.Second)
}
