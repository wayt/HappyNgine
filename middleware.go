package happy

import (
	"github.com/wayt/happyngine/context"
)

type MiddlewareHandler func(*context.Context) MiddlewareInterface

type MiddlewareInterface interface {
	HandleBefore() error
	HandleAfter() error
}

type Middleware struct {
	Context *context.Context
}
