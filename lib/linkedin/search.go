package linkedin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/goware/lg"
)

type Profile struct {
	FullName   string `json:"fullName"`
	ProfileURL string `json:"profile"`
}

// f_PC = company id (https://www.linkedin.com/company/2002577)
// orig = sorting, past company

func (c *Client) GetSearch(limit int) ([]*Profile, error) {
	urlx, _ := url.Parse("https://www.linkedin.com/vsearch/pj")
	params := url.Values{}
	params.Set("orig", "FCTD")
	params.Set("f_PC", fmt.Sprintf("%d", c.CompanyID))

	for i := 11; i <= limit; i++ {
		params.Set("page_num", fmt.Sprintf("%d", i))
		urlx.RawQuery = params.Encode()

		res, err := c.Get(urlx.String())
		if err != nil {
			lg.Fatal(err)
		}
		// TODO: json decode
		f, _ := os.Create(fmt.Sprintf("%s/%s_search_%d.json", DataDir, c.Company, i))
		io.Copy(f, res.Body)

		f.Close()
		res.Body.Close()

		lg.Printf("finished page %d of %d", i, limit)
		c.Sleep()
	}

	return nil, nil
}

func (c *Client) LoadParsedSearch() ([]*Profile, error) {
	f, err := os.Open(fmt.Sprintf("%s/%s.json", DataDir, c.Company))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var profiles []*Profile
	if err := json.NewDecoder(f).Decode(&profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}
