package happy

import (
    "net/http"
)

type Context struct{

    Request *http.Request
    Response http.ResponseWriter
    API *API
    Middlewares []MiddlewareInterface
}

func NewContext(req *http.Request, resp http.ResponseWriter, api *API) *Context {

    this := new(Context)

    this.Request = req
    this.Response = resp
    this.API = api

    return this
}
