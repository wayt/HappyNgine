package happy

import (
    "net/http"
)

type API struct {
}

type Context struct{

    Request *http.Request
    Response *http.Response
    API *API
}

func NewContext(req *http.Request, resp *http.Response, api *API) *Context {

    this := new(Context)

    this.Request = req
    this.Response = resp
    this.API = api

    return this
}
