package gplus

import (
	"encoding/json"
	"fmt"
	"github.com/wayt/happyngine/env"
	"io/ioutil"
	"net/http"
)

type TokenInfo struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Azp           string `json:"azp"`
	Email         string `json:"email"`
	AtHash        string `json:"at_hash"`
	EmailVerified string `json:"email_verified"`
	Aud           string `json:"aud"`
	Hd            string `json:"hd"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
}

func GetTokenInfo(token string) (*TokenInfo, error) {

	resp, err := http.Get(fmt.Sprintf(`https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s`, token))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	t := new(TokenInfo)
	if err := json.Unmarshal(data, t); err != nil {
		return nil, err
	}

	if env.Get("GOOGLEPLUS_CLIENT_ID") != t.Aud {
		return nil, nil
	}

	return t, nil
}
