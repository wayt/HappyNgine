package context

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
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

	c := new(Context)

	c.Request = req
	c.Response = resp
	c.UserData = make(map[string]interface{})
	c.ResponseStatusCode = 200

	c.RequestId = uuid.New()

	return this
}

func (c *Context) GetParam(key string) string {

	return c.Request.FormValue(key)
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

func (c *Context) GetURLParam(key string) string {

	return c.Request.URL.Query().Get(key)
}

func (c *Context) GetURLIntParam(key string) int {

	value, err := strconv.Atoi(c.GetURLParam(key))
	if err != nil {
		return 0
	}
	return value
}

func (c *Context) GetURLInt64Param(key string) int64 {

	value, err := strconv.ParseInt(c.GetURLParam(key), 10, 64)
	if err != nil {
		return 0
	}

	return value
}

func (c *Context) Send(code int, text string, headers ...string) {

	hasMime := false
	for _, header := range headers {

		array := strings.Split(header, ":")
		if len(array) != 2 {
			continue
		}

		c.Response.Header().Add(array[0], array[1])

		if array[0] == "Content-Type" {
			hasMime = true
		}
	}

	c.Response.Header().Add("X-happyngine-request-id", this.RequestId)
	if node := env.Get("NODE_NAME"); node != "" {
		c.Response.Header().Add("X-happyngine-node", node)
	}

	if !hasMime {
		c.Response.Header().Add("Content-Type", "application/json")
	}

	for k, v := range this.API.Headers {

		matchs := regexp.MustCompile(`^{(.*)}$`).FindStringSubmatch(v)
		if len(matchs) != 0 {
			header := matchs[1]
			if v = c.Request.Header.Get(header); len(v) == 0 {
				continue
			}
		}

		c.Response.Header().Add(k, v)
	}

	c.Response.WriteHeader(code)
	c.Response.Write([]byte(text))
	c.ResponseStatusCode = code
}

func (c *Context) RemoteIP() string {

	ipStr := strings.SplitN(c.Request.RemoteAddr, ":", 1)[0]

	if header := c.Request.Header.Get("X-Forwarded-For"); len(header) != 0 {
		ipStr = header
	}

	return ipStr
}

func (c *Context) Debugln(args ...interface{}) {
	debug.Println(append([]interface{}{c.RequestId}, args...))
}

func (c *Context) Warningln(args ...interface{}) {
	warning.Println(append([]interface{}{c.RequestId}, args...))
}

func (c *Context) Errorln(args ...interface{}) {
	err.Println(append([]interface{}{c.RequestId}, args...))
}

func (c *Context) Criticalln(args ...interface{}) {
	critical.Println(append([]interface{}{c.RequestId}, args...))
}
