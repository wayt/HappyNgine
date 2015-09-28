package happyngine

import (
	"net/http"
)

type ErrorHandler func(*Context, interface{})

type API struct {
	Middlewares     []MiddlewareHandler
	Resources       map[string]interface{}
	Error404Handler ErrorHandler
	Router          Router
	Headers         map[string]string
}

func NewAPI() *API {

	this := new(API)
	this.Resources = make(map[string]interface{})
	this.Error404Handler = this.error404Handler
	this.Headers = make(map[string]string)

	return this
}

type M map[string]interface{}

func (this *API) AddResource(name string, resource interface{}) {

	this.Resources[name] = resource
}

func (this *API) GetResource(name string) interface{} {

	if val, ok := this.Resources[name]; ok {
		return val
	}
	return nil
}

func (this *API) AddRoute(method string, path string, actionHandler ActionHandler, middlewares ...MiddlewareHandler) {

	this.Router.AddRoute(NewRoute(method, path, actionHandler, middlewares...))
}

func (this *API) AddMiddleware(middlewareHandler MiddlewareHandler) {

	this.Middlewares = append(this.Middlewares, middlewareHandler)
}

func (this *API) error404Handler(context *Context, err interface{}) {

	context.Response.WriteHeader(404)
	context.Response.Write([]byte("404 Not Found"))
}

func (this *API) preDispatch(route *Route, context *Context) error {

	context.middlewares = append(this.Middlewares, route.Middlewares...)
	context.action = route.ActionHandler

	return nil
}

func (this *API) dispatch(route *Route, c *Context) {

	c.Next()
}

func (this *API) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	context := NewContext(req, resp, this)

	route, err := this.Router.FindRoute(req)
	if err != nil {

		this.Error404Handler(context, err)
		return
	}

	if err := this.preDispatch(route, context); err != nil {
		return
	}

	this.dispatch(route, context)
}

func (this *API) Run(host string) error {

	return http.ListenAndServe(host, this)
}
