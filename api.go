package happy

import (
    "fmt"
)

type API struct {
    Resources map[string]interface{}
    Routes []*Route
}

func NewAPI() *API {

    this := new(API)
    this.Resources = make(map[string]interface{})

    return this
}

func(this *API) AddResource(name string, resource interface{}) {

    this.Resources[name] = resource
}

func(this *API) GetResource(name string) interface{} {

    return this.Resources[name]
}

func(this *API) AddRoute(method string, path string, actionHandler ActionHandler, middlewares ...MiddlewareHandler) {

    this.Routes = append(this.Routes, NewRoute(method, path, actionHandler, middlewares))
}

func(this *API) Run() {

    fmt.Println("Let's goooo")
}
