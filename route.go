package happy

type Route struct {
    Method string
    Path string
    ActionHandler ActionHandler
    Middlewares []MiddlewareHandler
}

func NewRoute(method string, path string, actionHandler ActionHandler, middlewares []MiddlewareHandler) *Route {

    this := new(Route)
    this.Method = method
    this.Path = path
    this.ActionHandler = actionHandler
    this.Middlewares = middlewares

    return this
}
