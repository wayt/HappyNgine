package happy

import (
    "net/http"
    "fmt"
)

type ErrorHandler func (*Context)

type API struct {
    Middlewares []MiddlewareHandler
    Resources map[string]interface{}
    ErrorHandler ErrorHandler
    Router Router
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

    this.Router.AddRoute(NewRoute(method, path, actionHandler, middlewares...))
}

func (this *API) AddMiddleware(middlewareHandler MiddlewareHandler) {

    this.Middlewares = append(this.Middlewares, middlewareHandler)
}

func (this *API) errorHandler(context *Context) {

    context.Response.WriteHeader(404)
    context.Response.Write([]byte("404 Not Found"))
}

func (this *API) preDispatch(route *Route, context *Context) error {

    for _, middlewareHandler := range append(this.Middlewares, route.Middlewares...) {

        m := middlewareHandler(context)
        context.Middlewares = append(context.Middlewares, m)
        if err := m.HandleBefore(); err != nil {

            return err
        }
    }

    return nil
}

func (this *API) dispatch(route *Route, context *Context) {

    action := route.ActionHandler(context)
    action.Run()
}

func (this *API) postDispatch(context *Context) error {

    for _, m := range context.Middlewares {

        if err := m.HandleAfter(); err != nil {

            return err
        }
    }

    return nil
}


func (this *API) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    fmt.Println(req.Method, ":", req.URL)

    context := NewContext(req, resp, this)

    route, err := this.Router.FindRoute(req)
    if err != nil {

        this.ErrorHandler(context)
        return
    }

    if err := this.preDispatch(route, context); err != nil {
        return
    }

    this.dispatch(route, context)

    if err := this.postDispatch(context); err != nil {
        return
    }
}

func (this *API) Run(host string) error {

    return http.ListenAndServe(host, this)
}
