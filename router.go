package happy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Router struct {
	Routes []*Route
}

func (this *Router) AddRoute(route *Route) {

	this.Routes = append(this.Routes, route)
}

func getMimeType(req *http.Request) string {

	return req.Header.Get("Content-Type")
}

func parseJsonBody(req *http.Request) {

	body, _ := ioutil.ReadAll(req.Body)

	bodyMap := make(map[string]interface{})
	json.Unmarshal(body, &bodyMap)

	for k, v := range bodyMap {

		req.Form.Add(k, fmt.Sprintf("%v", v))
	}

}

func (this *Router) parseRequestParams(req *http.Request, route *Route, matches []string) {

	// This will create the req.Form object
	req.ParseForm()

	mime := getMimeType(req)

	// Handle json input data
	if strings.Index(mime, "application/json") != -1 {

		parseJsonBody(req)
	}

	if len(matches) > 0 && matches[0] == req.URL.Path {

		for i, name := range route.Path.SubexpNames() {

			if len(name) > 0 {

				req.Form.Add(name, matches[i])
			}
		}
	}
}

func (this *Router) FindRoute(req *http.Request) (*Route, error) {

	for _, r := range this.Routes {

		if r.Method == req.Method {

			matches := r.Path.FindStringSubmatch(req.URL.Path)
			if matches != nil {

				this.parseRequestParams(req, r, matches)

				return r, nil
			}
		}
	}

	return nil, errors.New("No route")
}
