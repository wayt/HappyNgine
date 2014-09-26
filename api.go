package happy

import (
    "fmt"
    "net/http"
    "errors"
)

type API struct {
    Middlewares []MiddlewareHandler
    Resources map[string]interface{}
    Routes []*Route
}

func NewAPI() *API {

    this := new(API)
    this.Resources = make(map[string]interface{})

    return this
}

func (this *API) AddResource(name string, resource interface{}) {

    this.Resources[name] = resource
}

func (this *API) GetResource(name string) interface{} {

    return this.Resources[name]
}

func (this *API) AddRoute(method string, path string, actionHandler ActionHandler, middlewares ...MiddlewareHandler) {

    this.Routes = append(this.Routes, NewRoute(method, path, actionHandler, middlewares))
}

func (this *API) AddMiddleware(middlewareHandler MiddlewareHandler) {

    this.Middlewares = append(this.Middlewares, middlewareHandler)
}

func (this *API) findRouteForRequest(req *http.Request) (*Route, error) {

    for _, r := range this.Routes {

        if r.Path == req.URL.Path && r.Method == req.Method {

            return r, nil
        }
    }

    return nil, errors.New("No route")
}

func (this *API) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    fmt.Println(req.Method, ":", req.URL)

    context := NewContext(req, resp, this)

    route, err := this.findRouteForRequest(req)
    if err != nil {

        resp.WriteHeader(404)
        resp.Write([]byte("Not found"))
        return
    }

    var middlewares []MiddlewareInterface

    for _, middlewareHandler := range append(this.Middlewares, route.Middlewares...) {

        m := middlewareHandler(context)
        middlewares = append(middlewares, m)
        if err := m.HandleBefore(); err != nil {

            return
        }
    }

    action := route.ActionHandler(context)

    if action.IsValid() {

        action.Run()
    }

    for _, m := range middlewares {

        if err := m.HandleAfter(); err != nil {

            return
        }
    }
}

func (this *API) Run(host string) {

    fmt.Println("Let's goooo")

    http.ListenAndServe(host, this)
}
