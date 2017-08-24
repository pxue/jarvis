package mls

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Beanfield struct {
	Buildings []*Building `json:"buildings"`
}

type Building struct {
	Name    string `json:"name"`
	Status  string `"status"`
	IsOnNet bool   `json:"isOnnet"`
}

var BeanfieldApi = "https://api.beanfield.com/buildings"

func CheckBeanfield(postal string) (bool, error) {
	urlx, _ := url.Parse(BeanfieldApi)
	q := url.Values{
		"q": {postal},
	}
	urlx.RawQuery = q.Encode()

	res, err := http.Get(urlx.String())
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var bf *Beanfield
	if err := json.NewDecoder(res.Body).Decode(&bf); err != nil {
		return false, err
	}

	if len(bf.Buildings) == 0 {
		return false, errors.New("no building found")
	}

	return bf.Buildings[0].IsOnNet, nil
}
