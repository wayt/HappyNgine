package happy

type ActionHandler func(*Context) ActionInterface

type ActionInterface interface {

    IsValid() bool
    Run()
}

type Action struct {

    Context *Context
}
