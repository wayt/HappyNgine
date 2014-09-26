package happy

import (
    "net/http"
    "net/url"
    "errors"
)

type Router struct {
    Routes []*Route
}

func (this *Router) AddRoute(route *Route) {

    this.Routes = append(this.Routes, route)
}

func (this *Router) addRequestParam(req *http.Request, route *Route, matches []string) {

    if len(matches) > 0 && matches[0] == req.URL.Path {

        if req.PostForm == nil {

            req.PostForm = url.Values{}
        }

        for i, name := range route.Path.SubexpNames() {

            if len(name) > 0 {

                req.PostForm.Add(name, matches[i])
            }
        }
    }
}

func (this *Router) FindRoute(req *http.Request) (*Route, error) {

    for _, r := range this.Routes {

        if r.Method == req.Method {

            matches := r.Path.FindStringSubmatch(req.URL.Path)
            if matches != nil {

                this.addRequestParam(req, r, matches)

                return r, nil
            }
        }
    }

    return nil, errors.New("No route")
}
