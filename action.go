package happy

import (
    "github.com/gohappy/happy/validator"
)

type ActionHandler func(*Context) ActionInterface

type ActionInterface interface {

    Run()
}

type Action struct {

    Context *Context
    Parameters []*Parameter
    Errors []string
}

func (this *Action) IsValid() bool {

    request := this.Context.Request
    isValid := true

    for _, parameter := range this.Parameters {

        name := parameter.Name
        err := parameter.IsValid(request.FormValue(name))
        if err != nil {

            this.Errors = append(this.Errors, err.Error())
            isValid = false
        }
    }

   return isValid
}

func (this *Action) AddParameter(name string, validators ...*validator.Validator) {

    this.Parameters = append(this.Parameters, NewParameter(name, validators...))
}
