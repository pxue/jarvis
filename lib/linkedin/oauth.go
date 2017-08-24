package linkedin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/goware/lg"
	"github.com/pxue/jarvis/lib/ws"
)

type AuthResp struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

const (
	OAuthURL = "https://www.linkedin.com/uas/oauth2/authorization"
)

var (
	readScope = []string{
		"r_basicprofile",
	}
)

func OAuth(w http.ResponseWriter, r *http.Request) {
	scope := readScope

	args := url.Values{}
	args.Set("client_id", os.Getenv("CHURNER_APPID"))
	args.Set("response_type", "code")
	args.Set("redirect_uri", os.Getenv("CHURNER_REDIRECT"))
	args.Set("scope", strings.Join(scope, ","))
	args.Set("state", "churner_state_t0k3n")
	lg.Warnf("%+v", args)

	authUrl, _ := url.Parse(OAuthURL)
	authUrl.RawQuery = args.Encode()

	http.Redirect(w, r, authUrl.String(), 302)
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	cbArgs := r.URL.Query()
	code := cbArgs.Get("code")
	cbError := cbArgs.Get("error")
	cbErrorDesc := cbArgs.Get("error_description")

	// check state?
	if cbError != "" {
		ws.Respond(w, http.StatusUnauthorized, errors.New(fmt.Sprintf("%v: %v", cbError, cbErrorDesc)))
		return
	}

	// exchange code
	args := url.Values{}
	args.Add("client_id", os.Getenv("CHURNER_APPID"))
	args.Add("client_secret", os.Getenv("CHURNER_APPSECRET"))
	args.Add("redirect_uri", os.Getenv("CHURNER_REDIRECT"))
	args.Add("grant_type", "authorization_code")
	args.Add("code", code)
	lg.Warnf("%+v", args)

	resp, err := http.PostForm("https://www.linkedin.com/uas/oauth2/accessToken", args)
	if err != nil {
		ws.Respond(w, http.StatusUnauthorized, err)
		return
	}
	defer resp.Body.Close()

	var payload AuthResp
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		ws.Respond(w, http.StatusUnauthorized, err)
		return
	}

	ctx := context.WithValue(r.Context(), "access", &payload)
	GetProfile(w, r.WithContext(ctx))
}
