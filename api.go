package happy

import (
	"net/http"
	"strings"
)

type ErrorHandler func(*Context, interface{})

type API struct {
	Middlewares     []MiddlewareHandler
	Resources       map[string]interface{}
	Error404Handler ErrorHandler
	PanicHandler    ErrorHandler
	Router          Router
}

func NewAPI() *API {

	this := new(API)
	this.Resources = make(map[string]interface{})
	this.Error404Handler = this.error404Handler
	this.PanicHandler = this.panicHandler

	return this
}

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

func (this *API) panicHandler(context *Context, err interface{}) {

	context.Response.WriteHeader(500)
	context.Response.Write([]byte("Internal Server Error"))
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

	if action.IsValid() {

		action.Run()
	}

	errors, code := action.GetErrors()
	if len(errors) != 0 {

		response := `{"error":["` + strings.Join(errors, `","`) + `"]}`
		action.Send(code, response)
	}
}

func (this *API) postDispatch(context *Context) {

	for _, m := range context.Middlewares {

		m.HandleAfter()
	}
}

func (this *API) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	context := NewContext(req, resp, this)

	// Panic handler
	defer func() {
		if r := recover(); r != nil {

			this.PanicHandler(context, r)
		}
	}()

	route, err := this.Router.FindRoute(req)
	if err != nil {

		this.Error404Handler(context, err)
		return
	}

	if err := this.preDispatch(route, context); err != nil {

		// Excute HandleAfter for executed middleware
		this.postDispatch(context)
		return
	}

	this.dispatch(route, context)

	this.postDispatch(context)
}

func (this *API) Run() {

	http.HandleFunc("/", this.ServeHTTP)
}
