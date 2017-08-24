package linkedin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/goware/lg"
	"github.com/pxue/jarvis/lib/ws"
)

type ConnectedWrapper struct {
	Values []*Connected `json:"values"`
	Paging struct {
		Total int `json:"total"`
		Start int `json:"start"`
		Count int `json:"count"`
	} `json:"paging"`
}

type Connected struct {
	Addresses []interface{} `json:"addresses"`
	Company   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"company"`
	ConnectionDate        int64         `json:"connectionDate"`
	Created               int64         `json:"created"`
	DisplaySources        []string      `json:"displaySources"`
	Emails                []interface{} `json:"emails"`
	FamiliarName          string        `json:"familiarName"`
	FirstName             string        `json:"firstName"`
	FullName              string        `json:"fullName"`
	GraphDistance         string        `json:"graphDistance"`
	LastName              string        `json:"lastName"`
	LegacyMergedContactID string        `json:"legacyMergedContactId"`
	Location              string        `json:"location"`
	MemberID              int64         `json:"memberId"`
	MergedContactID       string        `json:"mergedContactId"`
	Name                  string        `json:"name"`
	PhoneNumbers          []struct {
		Number     string `json:"number"`
		Primary    bool   `json:"primary"`
		SourceType string `json:"sourceType"`
		Type       string `json:"type"`
		Visible    bool   `json:"visible"`
	} `json:"phoneNumbers"`
	ProfileImageURL string   `json:"profileImageUrl"`
	ProfileURL      string   `json:"profileUrl"`
	Sources         []string `json:"sources"`
	Tags            []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
	Title string `json:"title"`
}

const (
	PerPage = 20
)

func queryLinkedin(uri, accessToken string, args url.Values) (b []byte, statusCode int, err error) {
	if args == nil {
		args = url.Values{}
	}
	args.Set("format", "json")

	u, _ := url.Parse(ApiUrl + uri)
	u.RawQuery = args.Encode()
	uri = u.String()
	lg.Warn(uri)

	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	if err != nil {
		lg.Warningf("querying linkedin uri:%v failed:%v", uri, err.Error())
		return nil, 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		lg.Warningf("getting response from linkedin uri:%v failed:%v", uri, err.Error())
		return nil, 0, err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		lg.Warningf("reading response for linkedin uri%v failed:%v", err.Error())
		return nil, 0, err

	}
	return b, resp.StatusCode, nil
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	access := r.Context().Value("access").(*AuthResp)
	lg.Warnf("%+v", access)

	queryUri := "/v1/people/~:(id,positions:(id,title,summary,start-date,end-date,is-current))"
	b, _, _ := queryLinkedin(queryUri, access.AccessToken, nil)
	ws.Respond(w, http.StatusOK, string(b))
}

func linkedinCall(reqUrl *url.URL, params url.Values) (*ConnectedWrapper, error) {
	rawCookies := ""
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", rawCookies)
	req, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))

	req.URL = reqUrl
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	req.RequestURI = ""

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var payload *ConnectedWrapper
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func GetConnected() error {
	params, _ := url.ParseQuery("fields=id,name,firstName,lastName,company,title,location,tags,emails,sources,displaySources,connectionDate,secureProfileImageUrl&sort=CREATED_DESC&_=1482966483025")
	params.Set("start", "80")
	params.Set("count", fmt.Sprintf("%d", PerPage))

	reqUrl, _ := url.Parse("https://www.linkedin.com/connected/api/v2/contacts")
	initCall, err := linkedinCall(reqUrl, params)
	if err != nil {
		return err
	}

	var connections []*Connected
	for p := initCall.Paging.Start; p < initCall.Paging.Total; p += PerPage {
		params.Set("start", fmt.Sprintf("%d", p))
		params.Set("count", fmt.Sprintf("%d", PerPage))
		payload, err := linkedinCall(reqUrl, params)
		if err != nil {
			lg.Error(err)
			continue
		}

		if len(payload.Values) == 0 {
			break
		}
		connections = append(connections, payload.Values...)
		lg.Warnf("parsed: start %d of total %d", p, initCall.Paging.Total)
		break
	}

	cnFile, err := os.Create("conn.json")
	if err != nil {
		return err
	}
	defer cnFile.Close()

	return json.NewEncoder(cnFile).Encode(connections)
}

func (c *Client) LoadConnected() ([]*Connected, error) {
	f, err := os.Open("conn.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var connections []*Connected
	if err := json.NewDecoder(f).Decode(&connections); err != nil {
		return nil, err
	}

	return connections, nil
}
