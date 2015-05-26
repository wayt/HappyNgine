package happy

import (
	"github.com/gohappy/happy/context"
)

type MiddlewareHandler func(*context.Context) MiddlewareInterface

type MiddlewareInterface interface {
	HandleBefore() error
	HandleAfter() error
}

type Middleware struct {
	Context *context.Context
}
