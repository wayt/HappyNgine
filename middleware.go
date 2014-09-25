package happy

import (
    "fmt"
)

type Middleware struct{

    Ctx *Context
}

type MiddlewareHandler func(*Context) *Middleware

func NewMiddleware(ctx *Context) *Middleware {

    this := new(Middleware)

    this.Ctx = ctx

    return this
}

func (this Middleware) HandleBefore() error {

    fmt.Println("Middleware.HandleBefore")
    return nil
}

func (this Middleware) HandleAfter() error {

    fmt.Println("Middleware.HandleAfter")
    return nil
}

