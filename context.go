package happy

import (
    "net/http"
)

type Context struct{

    Request *http.Request
    Response http.ResponseWriter
    API *API
    Middlewares []MiddlewareInterface
    UserData map[string]map[string]interface{}
    ResponseStatusCode int // Because we can't retrieve the status from http.ResponseWriter
}

func NewContext(req *http.Request, resp http.ResponseWriter, api *API) *Context {

    this := new(Context)

    this.Request = req
    this.Response = resp
    this.API = api
    this.UserData = make(map[string]map[string]interface{})
    this.ResponseStatusCode = 200

    return this
}

func (this *Context) GetParam(key string) string {

    return this.Request.FormValue(key)
}

func (this *Context) Send(code int, text string) {

        this.Response.WriteHeader(code)
        this.Response.Write([]byte(text))
        this.ResponseStatusCode = code
}
