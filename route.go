package happy

import (
	"fmt"
	"regexp"
)

type Route struct {
	Method        string
	Path          *regexp.Regexp
	ActionHandler ActionHandler
	Middlewares   []MiddlewareHandler
}

func NewRoute(method string, path string, actionHandler ActionHandler, middlewares ...MiddlewareHandler) *Route {

	this := new(Route)
	this.Method = method

	r := regexp.MustCompile(`:[^/#?()\.\\]+`)
	path = r.ReplaceAllStringFunc(path, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})

	var err error
	if this.Path, err = regexp.Compile("^" + path + "$"); err != nil {
		panic(err)
	}
	this.ActionHandler = actionHandler
	this.Middlewares = middlewares

	return this
}
