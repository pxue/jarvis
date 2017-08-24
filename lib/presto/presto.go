package presto

var (
	RootURL    = "https://www.prestocard.ca"
	SignInPath = "/api/sitecore/AFMSAuthentication/SignInWithAccount"
)

func Signin(username, password string) {

	var payload struct {
		CustSecurity struct {
			Login    string
			Password string
		} `json:"CustSecurity"`
	}
	_ = payload

}
