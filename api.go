package happy

import (
    "net/http"
)

type ErrorHandler func (*Context)

type API struct {
    Middlewares []MiddlewareHandler
    Resources map[string]interface{}
    Routes []*Route
    ErrorHandler ErrorHandler
}

func NewAPI() *API {

    this := new(API)
    this.Resources = make(map[string]interface{})
    this.ErrorHandler = this.errorHandler

    return this
}

func (this *API) AddResource(name string, resource interface{}) {

    this.Resources[name] = resource
}

func (this *API) GetResource(name string) interface{} {

    return this.Resources[name]
}

func (this *API) AddRoute(method string, path string, actionHandler ActionHandler, middlewares ...MiddlewareHandler) {

    this.Routes = append(this.Routes, NewRoute(method, path, actionHandler, middlewares...))
}

func (this *API) AddMiddleware(middlewareHandler MiddlewareHandler) {

    this.Middlewares = append(this.Middlewares, middlewareHandler)
}

func (this *API) errorHandler(context *Context) {

    context.Response.WriteHeader(404)
    context.Response.Write([]byte("404 Not Found"))
}

func (this *API) findRouteForRequest(req *http.Request) *Route {

    for _, route := range this.Routes {

        if route.Path == req.URL.Path && route.Method == req.Method {

            return route
        }
    }

    return nil
}

func (this *API) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

    context := NewContext(req, resp, this)

    route := this.findRouteForRequest(req)

    // If route not found
    if route == nil {

        this.ErrorHandler(context)
        return
    }

    var middlewares []*Middleware

    // Predispatch
    for _, middlewareHandler := range append(this.Middlewares, route.Middlewares...) {

        m := middlewareHandler(context)
        middlewares = append(middlewares, m)
        if err := m.HandleBefore(); err != nil {

            return
        }
    }

    // Do an action
    action := route.ActionHandler(context)
    action.Run()

    // Postdispatch
    for _, m := range middlewares {

        if err := m.HandleAfter(); err != nil {

            return
        }
    }
}

func (this *API) Run(host string) {

    http.ListenAndServe(host, this)
}
