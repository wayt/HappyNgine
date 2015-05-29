package happyngine

import (
	"github.com/wayt/happyngine/validator"
)

type ActionHandler func(*Context) ActionInterface

type ActionInterface interface {
	Run()
	IsValid() bool
	GetErrors() ([]string, int)
	Send(int, string, ...string)
	HasErrors() bool
}

type Action struct {
	Context    *Context
	Parameters []*Parameter
	Errors     []string
	ErrorCode  int
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

func (this *Action) HasErrors() bool {
	return len(this.Errors) != 0
}

func (this *Action) GetErrors() ([]string, int) {

	return this.Errors, this.ErrorCode
}

func (this *Action) AddParameter(name string, required bool, validators ...*validator.Validator) {

	this.Parameters = append(this.Parameters, NewParameter(name, required, validators...))
}

func (this *Action) AddError(code int, text string) {

	this.ErrorCode = code
	this.Errors = append(this.Errors, text)
}

func (this *Action) GetParam(key string) string {

	return this.Context.GetParam(key)
}

func (this *Action) GetIntParam(key string) int {

	return this.Context.GetIntParam(key)
}

func (this *Action) GetInt64Param(key string) int64 {

	return this.Context.GetInt64Param(key)
}

func (this *Action) GetURLParam(key string) string {

	return this.Context.GetURLParam(key)
}

func (this *Action) GetURLIntParam(key string) int {

	return this.Context.GetURLIntParam(key)
}

func (this *Action) GetURLInt64Param(key string) int64 {

	return this.Context.GetURLInt64Param(key)
}

func (this *Action) Send(code int, text string, headers ...string) {

	this.Context.Send(code, text, headers...)
}
