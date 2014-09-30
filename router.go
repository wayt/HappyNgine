package happy

import (
    "net/http"
    "errors"
)

type Router struct {
    Routes []*Route
}

func (this *Router) AddRoute(route *Route) {

    this.Routes = append(this.Routes, route)
}

func (this *Router) addRequestParam(req *http.Request, route *Route, matches []string) {

    req.ParseForm()

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

                this.addRequestParam(req, r, matches)

                return r, nil
            }
        }
    }

    return nil, errors.New("No route")
}
