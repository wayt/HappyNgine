package context

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/gohappy/happy/env"
	"regexp"
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
	RequestId          string
}

func NewContext(req *http.Request, resp http.ResponseWriter) *Context {

	this := new(Context)

	this.Request = req
	this.Response = resp
	this.UserData = make(map[string]interface{})
	this.ResponseStatusCode = 200

	this.RequestId = uuid.New()

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

func (this *Context) GetURLParam(key string) string {

	return this.Request.URL.Query().Get(key)
}

func (this *Context) GetURLIntParam(key string) int {

	value, err := strconv.Atoi(this.GetURLParam(key))
	if err != nil {
		return 0
	}
	return value
}

func (this *Context) GetURLInt64Param(key string) int64 {

	value, err := strconv.ParseInt(this.GetURLParam(key), 10, 64)
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

	this.Response.Header().Add("X-happyngine-request-id", this.RequestId)
	if node := env.Get("NODE_NAME"); node != "" {
		this.Response.Header().Add("X-happyngine-node", node)
	}

	if !hasMime {
		this.Response.Header().Add("Content-Type", "application/json")
	}

	for k, v := range this.API.Headers {

		matchs := regexp.MustCompile(`^{(.*)}$`).FindStringSubmatch(v)
		if len(matchs) != 0 {
			header := matchs[1]
			if v = this.Request.Header.Get(header); len(v) == 0 {
				continue
			}
		}

		this.Response.Header().Add(k, v)
	}

	this.Response.WriteHeader(code)
	this.Response.Write([]byte(text))
	this.ResponseStatusCode = code
}

func (this *Context) RemoteIP() string {

	ipStr := strings.SplitN(this.Request.RemoteAddr, ":", 1)[0]

	if header := this.Request.Header.Get("X-Forwarded-For"); len(header) != 0 {
		ipStr = header
	}

	return ipStr
}
