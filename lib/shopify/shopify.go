package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/goware/lg"
)

func worker(sg *sync.WaitGroup, site chan string) {
	for st := range site {
		resp, err := http.Head(st)
		if err != nil {
			continue
		}

		var isShopify bool
		for _, v := range resp.Header {
			if !strings.Contains(v[0], "Shopify") {
				continue
			}

			isShopify = true
			break
		}

		if !isShopify {
			continue
		}

		resp, err = http.Get(st)
		if err != nil {
			continue
		}

		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			continue
		}

		var salesLink string
		doc.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			link, _ := s.Attr("href")
			linkText := strings.ToLower(strings.TrimSpace(s.Text()))
			if strings.Contains(link, "collections") &&
				strings.Contains(linkText, "sale") {
				salesLink = st + link
				return false
			}
			return true
		})

		if salesLink == "" {
			lg.Warn("no sale link")
			continue
		}

		lg.Warn(salesLink)
		resp, err = http.Get(salesLink + ".oembed")
		if err != nil {
			continue
		}

		u, err := url.Parse(st)
		if err != nil {
			lg.Warnf("errored %s with %v", st, err)
			continue
		}

		f, err := os.Create(fmt.Sprintf("./data/%s.json", u.Host))
		if err != nil {
			lg.Warnf("errored %s with %v", st, err)
			continue
		}

		ioutil.ReadAll(resp.Body)
		io.Copy(f, resp.Body)

		f.Close()
		resp.Body.Close()

	}
	sg.Done()
}

func getProducts() {

	//locale, err := data.DB.Locale.FindOne(db.Cond{"shorthand": "queen-west"})
	//if err != nil {
	//lg.Fatal(err)
	//}

	sitesCh := make(chan string)
	var sg sync.WaitGroup
	sg.Add(10)
	for i := 1; i <= 10; i++ {
		go worker(&sg, sitesCh)
	}

	//result := data.DB.Place.Find(db.Cond{"locale_id": locale.ID})

	sites := []string{

	// sales
	//"http://thestoreonqueen.com/",
	//"https://duewest.ca",
	//"http://shop.thelegendsleague.com/",
	//"http://hayleyelsaesser.com/",
	//"http://smoke-ash.myshopify.com/",
	//"https://fasinfrankvintage.com/",
	//"http://shop.exclucitylife.com",
	//"http://www.livify.ca/",

	// no sales
	//"http://runwayluxe.com/",
	//"https://www.yo-sox.ca/",
	//"http://www.nobis.ca/",
	//"http://untitledandco.com/",
	//"http://www.ekojewellery.com",
	//"http://shop.getfreshcompany.com/",
	//"http://www.baileynelson.ca/",
	//"http://nuvango.com/",
	//"http://thecureapothecary.com/",
	//"http://www.coalminersdaughter.ca/",
	//"http://www.blurmakeuproom.com/",
	//"http://www.titikaactive.com/",
	}

	for _, s := range sites {
		//var pl *data.Place
		//if !result.Next(&pl) {
		//break
		//}

		//if pl.Website == "" {
		//continue
		//}

		sitesCh <- s
	}
	//result.Close()
	close(sitesCh)

	sg.Wait()
}

type Offer struct {
	CurrencyCode string  `json:"currency_code"`
	InStock      bool    `json:"in_stock"`
	OfferID      int64   `json:"offer_id"`
	Price        float64 `json:"price"`
	Sku          string  `json:"sku"`
	Title        string  `json:"title"`
}

type Product struct {
	ProductID    string   `json:"product_id"`
	Brand        string   `json:"brand"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Offers       []*Offer `json:"offers"`
}

func parseProduct() {
	f, err := os.Open("./data/duewest.ca.json")
	if err != nil {
		lg.Fatal(err)
	}

	var wrapper struct {
		Products []*Product `json:"products"`
		Provider string     `json:"provider"`
	}
	if err := json.NewDecoder(f).Decode(&wrapper); err != nil {
		lg.Fatal(err)
	}

	for _, p := range wrapper.Products {
		lg.Info(p.Title)
		for _, o := range p.Offers {
			lg.Info(o.Price)
			break
		}
	}
}

func main() {
	//conf := &data.DBConf{
	//Database:        "localyyz",
	//Hosts:           []string{"localhost:5432"},
	//Username:        "localyyz",
	//ApplicationName: "promo loader",
	//}
	//if _, err := data.NewDBSession(conf); err != nil {
	//log.Fatalf("db err: %s", err)
	//}

	parseProduct()
}
