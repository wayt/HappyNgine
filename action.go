package happy

import (
    "github.com/gohappy/happy/validator"
)

type ActionHandler func(*Context) ActionInterface

type ActionInterface interface {

    Run()
    IsValid() bool
    GetErrors() ([]string, int)
    Send(int, string)
}

type Action struct {

    Context *Context
    Parameters []*Parameter
    Errors []string
    ErrorCode int
}

func (this *Action) IsValid() bool {

    request := this.Context.Request
    isValid := true

    for _, parameter := range this.Parameters {

        name := parameter.Name
        err := parameter.IsValid(request.FormValue(name))
        if err != nil {

            this.AddError(400, err.Error())
            isValid = false
        }
    }

   return isValid
}

func (this *Action) GetErrors() ([]string, int) {

    return this.Errors, this.ErrorCode
}

func (this *Action) AddParameter(name string, validators ...*validator.Validator) {

    this.Parameters = append(this.Parameters, NewParameter(name, validators...))
}

func (this *Action) AddError(code int, text string) {

    this.ErrorCode = code
    this.Errors = append(this.Errors, text)
}

func (this *Action) GetParam(key string) string {

    return this.Context.GetParam(key)
}

func (this *Action) Send(code int, text string) {

    this.Context.Send(code, text)
}
