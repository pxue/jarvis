package form

// inspired by github.com/henrylee2cn/pholcus crawler

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goware/lg"
)

// Form is the default form element.
type Form struct {
	selection *goquery.Selection
	method    string
	action    string
	fields    url.Values
	buttons   url.Values
	cookies   []*http.Cookie
}

// NewForm parses form attributes and data from a goquery
// form selection
func NewForm(parentUrl string, s *goquery.Selection) *Form {
	form := &Form{
		selection: s,
		method:    "GET",
		action:    parentUrl,
		fields:    url.Values{},
		buttons:   url.Values{},
	}
	form.parse()

	return form
}

func (f *Form) SetInput(name, value string) {
	f.fields.Set(name, value)
}

// submit the form
func (f *Form) Submit() (*http.Response, error) {
	formVal := url.Values{}
	for name, val := range f.fields {
		lg.Printf("field %s: %s", name, val)
		formVal[name] = val
	}

	// add the button value
	for name, val := range f.buttons {
		formVal[name] = val
	}

	// assume it's always post for now
	req, err := http.NewRequest("POST", f.action, bytes.NewBufferString(formVal.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return http.DefaultClient.Do(req)
}

// serializes the form
func (f *Form) parse() {
	if method, found := f.selection.Attr("method"); found && method != "" {
		f.method = strings.ToUpper(method)
	}
	if action, found := f.selection.Attr("action"); found && action != "" {
		f.action = action
	}

	// parse inputs
	inputs := f.selection.Find("input,button,textarea")
	if inputs.Length() == 0 {
		return
	}

	inputs.Each(func(_ int, s *goquery.Selection) {
		name, found := s.Attr("name")
		if !found {
			return
		}

		typ, found := s.Attr("type")
		if !found && !s.Is("textarea") {
			return
		}

		if typ == "submit" {
			val, _ := s.Attr("value")
			f.buttons.Add(name, val)
			return
		}

		val, _ := s.Attr("value")
		f.fields.Add(name, val)
	})
}
