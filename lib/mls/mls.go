package mls

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	db "upper.io/db.v2"

	"github.com/PuerkitoBio/goquery"
	"github.com/goware/lg"
	"github.com/pxue/jarvis/lib/gmaps"
	"github.com/pxue/jarvis/lib/mls/data"
)

const (
	// table header column positions
	AddressHdr = 1
	AptHdr     = 2
	PriceHdr   = 4
	MLSHdr     = 8
)

const (
	Work = "488 wellington st w, toronto, canada"
)

func ParseListings() error {
	// maps client
	maps, err := gmaps.New()
	if err != nil {
		return err
	}

	DB, err := data.NewDBSession()
	if err != nil {
		return err
	}

	res, err := http.Get("http://v3.torontomls.net/Live/Pages/Public/Link.aspx?Key=f051e20fd6e043c484edb8e1e6551a69&App=TREB")
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return err
	}

	numListing := doc.Find(".data-list-row-number").Size()

	listingMap := make(map[string]*data.Listing, numListing)
	doc.Find("table tbody tr").Each(func(row int, tr *goquery.Selection) {
		l := &data.Listing{}
		tr.Find("td").Each(func(col int, td *goquery.Selection) {
			switch col {
			case AddressHdr:
				l.Address = td.Text()
			case AptHdr:
				l.Unit, _ = strconv.Atoi(td.Text())
			case PriceHdr:
				p := td.Text()[1:2] + td.Text()[3:6]
				l.Price, _ = strconv.Atoi(p)
			case MLSHdr:
				l.MLS = td.Text()
			default:
				// do nothing, skip
			}
		})

		// database call, check if already parsed
		dbListing, err := DB.Listing.FindByMLS(l.MLS)
		if err != nil {
			if err != db.ErrNoMoreRows {
				lg.Errorf("failed to query %s listing: %v", l.MLS, err)
			}
		}
		if dbListing != nil {
			listingMap[l.MLS] = dbListing
			return
		}

		listingMap[l.MLS] = l
	})

	doc.Find("div.reports div.link-item").Each(func(_ int, s *goquery.Selection) {
		MLS, _ := s.Attr("id")
		l := listingMap[MLS]
		if l.ID != 0 {
			//lg.Debugf("skipping %s", MLS)
			return
		}

		detLink := s.AttrOr(
			"data-deferred-loaded",
			s.AttrOr("data-deferred-load", ""),
		)
		if detLink == "" {
			lg.Warn("no defer load link found")
			return
		}
		detDoc, err := getDetail(detLink)
		if err != nil {
			lg.Warnf("detail doc: %v", err)
			return
		}

		reportDiv := "div[class*='status-'] "

		// apt size
		aptSizeLabel := detDoc.Find(reportDiv + "label:contains('Apx Sqft')")
		aptSizeRaw := aptSizeLabel.SiblingsFiltered("span.value").Text()
		l.Size.UnmarshalText(aptSizeRaw)

		// locker
		lockerLabel := detDoc.Find(reportDiv + "label:contains('Locker:')")
		l.HasLocker = (lockerLabel.SiblingsFiltered("span.value").Text() != "None")

		// walking distance
		l.Distance, err = maps.GetDistance(Work, fmt.Sprintf("%s, Toronto", l.Address))
		if err != nil {
			lg.Warn(err)
			return
		}

		// Beanfield Api
		// TODO: postal code
		l.HasBeanField, err = CheckBeanfield("")
		if err != nil {
			lg.Warn(err)
		}

		// save to db
		if err := DB.Save(l); err != nil {
			lg.Errorf("failed to save %+v with %v", l, err)
		}
		lg.Debugf("saved %s", MLS)
	})

	return nil
}

func getDetail(link string) (*goquery.Document, error) {
	linkx, _ := url.Parse(link)
	linkx.RawQuery = fmt.Sprintf("%s&_=%d", linkx.RawQuery, time.Now().UnixNano())

	res, err := http.Get(linkx.String())
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromResponse(res)
}
