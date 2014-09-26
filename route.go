package happy

import (
    "regexp"
)

type Route struct {
    Method string
    Path *regexp.Regexp
    ActionHandler ActionHandler
    Middlewares []MiddlewareHandler
}

func NewRoute(method string, path string, actionHandler ActionHandler, middlewares ...MiddlewareHandler) *Route {

    this := new(Route)
    this.Method = method
    var err error
    if this.Path, err = regexp.Compile(path); err != nil {
        panic(err)
    }
    this.ActionHandler = actionHandler
    this.Middlewares = middlewares

    return this
}
