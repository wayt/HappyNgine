package happyngine

import (
	//"code.google.com/p/go-uuid/uuid"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var hostname = ""

func init() {

	var err error
	hostname, err = os.Hostname()
	if err != nil || hostname == "" {
		hostname = env.Get("NODE_NAME")
	}
}

type Context struct {
	Request            *http.Request          `json:"request"`
	Response           http.ResponseWriter    `json:"-"`
	Middlewares        []MiddlewareInterface  `json:"-"`
	API                *API                   `json:"-"`
	UserData           map[string]interface{} `json:"user_data"`
	ResponseStatusCode int                    `json:"-"` // Because we can't retrieve the status from http.ResponseWriter
	RequestId          string                 `json:"request_id"`
	Errors             map[string]string      `json:"-"`
	ErrorCode          int                    `json:"-"`
}

func NewContext(req *http.Request, resp http.ResponseWriter, api *API) *Context {

	c := new(Context)

	c.Request = req
	c.Response = resp
	c.API = api
	c.UserData = make(map[string]interface{})
	c.ResponseStatusCode = 200
	c.Errors = make(map[string]string)

	c.RequestId = "42" // uuid.New()

	return c
}

func (c *Context) GetParam(key string) string {

	return c.Request.FormValue(key)
}

func (c *Context) GetIntParam(key string) int {

	value, err := strconv.Atoi(c.Request.FormValue(key))
	if err != nil {
		return 0
	}

	return value
}

func (c *Context) GetInt64Param(key string) int64 {

	value, err := strconv.ParseInt(c.Request.FormValue(key), 10, 64)
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
	c.SendByte(code, []byte(text), headers...)
}

func (c *Context) SendByte(code int, data []byte, headers ...string) {

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

	c.Response.Header().Add("X-happyngine-request-id", c.RequestId)
	if hostname != "" {
		c.Response.Header().Add("X-happyngine-node", hostname)
	}

	if !hasMime {
		c.Response.Header().Add("Content-Type", "application/json")
	}

	for k, v := range c.API.Headers {

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
	c.Response.Write(data)
	c.ResponseStatusCode = code
}

func (c *Context) RemoteIP() string {

	ipStr := strings.SplitN(c.Request.RemoteAddr, ":", 2)[0]

	if header := c.Request.Header.Get("X-Forwarded-For"); len(header) != 0 {
		// Because of google http load balancer
		// X-Forwarded-For: <client IP(s)>, <global forwarding rule external IP> (requests only)
		ipStr = strings.Split(header, ",")[0]
	}

	return strings.Trim(ipStr, " ")
}

func (c *Context) Debugln(args ...interface{}) {
	log.Debugln(append([]interface{}{c.RequestId}, args...))
}

func (c *Context) Warningln(args ...interface{}) {
	log.Warningln(append([]interface{}{c.RequestId}, args...))
}

func (c *Context) Errorln(args ...interface{}) {
	log.Errorln(append([]interface{}{c.RequestId}, args...))
}

func (c *Context) Criticalln(args ...interface{}) {
	log.Criticalln(append([]interface{}{c.RequestId}, args...))
}

func (c *Context) AddError(code int, text string) {
	c.ErrorCode = code
	c.Errors[text] = text
}

func (c *Context) GetErrors() ([]string, int) {

	errs := make([]string, 0)
	for _, err := range c.Errors {
		errs = append(errs, err)
	}

	return errs, c.ErrorCode
}

func (c *Context) HasErrors() bool {
	return len(c.Errors) != 0
}
