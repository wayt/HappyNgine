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

// access_token: ya29.EwI68PHHYxN48SvORyq-Y4XQ76GV6cj9VGZ5TW7wwqzibcTkiet1aON2MQCN2m_EEb1L
// id_token: eyJhbGciOiJSUzI1NiIsImtpZCI6IjdjYjg2MDAyODgyMTg0ZWVjMDVlMGFmY2U2NmY5ZmY4ZTA1YjE3MTMifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXRfaGFzaCI6InJWNEh5dy03aWpVLUtEX3hqbGYzbFEiLCJhdWQiOiI1Mzg4MjA1ODE1MDgtN2hwZjA2ZzNjZWdkc3BvMGlrb2hhNzR0cDFyaGhtM3AuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTAwNDAyOTE3Nzg4ODk0OTAzOTUiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXpwIjoiNTM4ODIwNTgxNTA4LTdocGYwNmczY2VnZHNwbzBpa29oYTc0dHAxcmhobTNwLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiaGQiOiJteS1zaWduLmNvbSIsImVtYWlsIjoibWF4aW1lQG15LXNpZ24uY29tIiwiaWF0IjoxNDQ1NDI4NTk3LCJleHAiOjE0NDU0MzIxOTcsIm5hbWUiOiJNYXhpbWUgR2ludGVycyIsInBpY3R1cmUiOiJodHRwczovL2xoNi5nb29nbGV1c2VyY29udGVudC5jb20vLWFmSE1TM1FZWVJBL0FBQUFBQUFBQUFJL0FBQUFBQUFBQUJBL0NNYlRDWTBqR1RJL3M5Ni1jL3Bob3RvLmpwZyIsImdpdmVuX25hbWUiOiJNYXhpbWUiLCJmYW1pbHlfbmFtZSI6IkdpbnRlcnMiLCJsb2NhbGUiOiJlbiJ9.aCT0UcqIxnIEEeeY6TE-pbqWla-UptGfZWSvwj2AHO403rbgNhrBhc4sze2DBMrh6pmSrxdqEY8Gj_arWtMAPRg8Gi3h302DzLtm_2yCRWYQgDXgM9Sz8Ay0fATRZeWTtBTJsmT5bkJoVivA8CxMOGh1x73Hs_mWIGpV9GXJP_jHVTx0aRwM1rGRioiQWBwQ_vCDOQe30otbGESwrawMNFZCtOZpwNuUiCUpmpZfVURV0N3vthoxomTiVU_PM8GE7_isg12AxZqz_A2IUSBBVajEVOqlkX-EySOfN0ZfJ6nSXIiVz0VatAKBqunq2fMB7Hi03pwlEMQ912N13nRe2w

func GetTokenInfo(token string, mobile bool) (*TokenInfo, error) {

	tokenType := "id_token"
	if mobile {
		tokenType = "access_token"
	}

	resp, err := http.Get(fmt.Sprintf(`https://www.googleapis.com/oauth2/v3/tokeninfo?%s=%s`, tokenType, token))
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
