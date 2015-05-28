package happyngine

type MiddlewareHandler func(*Context) MiddlewareInterface

type MiddlewareInterface interface {
	HandleBefore() error
	HandleAfter() error
}

type Middleware struct {
	Context *Context
}
