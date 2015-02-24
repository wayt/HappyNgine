package happy

import (
	"net/http"
	"strconv"
	"strings"
)

type Context struct {
	Request            *http.Request
	Response           http.ResponseWriter
	API                *API
	Middlewares        []MiddlewareInterface
	UserData           map[string]interface{}
	ResponseStatusCode int // Because we can't retrieve the status from http.ResponseWriter
}

func NewContext(req *http.Request, resp http.ResponseWriter, api *API) *Context {

	this := new(Context)

	this.Request = req
	this.Response = resp
	this.API = api
	this.UserData = make(map[string]interface{})
	this.ResponseStatusCode = 200

	return this
}

func (this *Context) GetParam(key string) string {

	return this.Request.FormValue(key)
}

func (this *Context) GetIntParam(key string) int {

	value, err := strconv.Atoi(this.Request.FormValue(key))
	if err != nil {
		return 0
	}

	return value
}

func (this *Context) GetInt64Param(key string) int64 {

	value, err := strconv.ParseInt(this.Request.FormValue(key), 10, 64)
	if err != nil {
		return 0
	}

	return value
}

func (this *Context) Send(code int, text string, headers ...string) {

	hasMime := false
	for _, header := range headers {

		array := strings.Split(header, ":")
		if len(array) != 2 {
			continue
		}

		this.Response.Header().Add(array[0], array[1])

		if array[0] == "Content-Type" {
			hasMime = true
		}
	}

	if !hasMime {
		this.Response.Header().Add("Content-Type", "application/json")
	}

	if this.API.AllowOrigin != "" {
		this.Response.Header().Add("Access-Control-Allow-Origin", this.API.AllowOrigin)
	}

	this.Response.WriteHeader(code)
	this.Response.Write([]byte(text))
	this.ResponseStatusCode = code
}
