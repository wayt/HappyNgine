package happy

type ActionHandler func(*Context) *Action

type Action struct {

    Context *Context

}

func (this *Action) IsValid() bool {

    return true
}

func (this *Action) Run() {

}
