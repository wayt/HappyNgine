package happyngine

import (
	"github.com/wayt/happyngine/validator"
)

type ActionHandler func(*Context) ActionInterface

type ActionInterface interface {
	Run()
	IsValid() bool
	Send(int, string, ...string)
}

type Action struct {
	Context    *Context
	Form       *Form
	Parameters []*Parameter
}

func (this *Action) IsValid() bool {

	if !this.oldIsValid() {
		return false
	}

	// In case of custom AddError in action New
	if this.HasErrors() {
		return false
	}

	if this.Form == nil {
		return true
	}

	return this.Form.IsValid()
}

func (this *Action) oldIsValid() bool {

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

	return this.Context.HasErrors()
}

func (this *Action) GetErrors() ([]string, int) {

	return this.Context.GetErrors()
}

func (this *Action) AddParameter(name string, required bool, validators ...*validator.Validator) {

	this.Parameters = append(this.Parameters, NewParameter(name, required, validators...))
}

func (this *Action) AddError(code int, text string) {

	this.Context.AddError(code, text)
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

	this.Context.SendByte(code, []byte(text), headers...)
}

func (this *Action) SendByte(code int, data []byte, headers ...string) {

	this.Context.SendByte(code, data, headers...)
}

func (this *Action) JSON(code int, obj interface{}) {
	this.Context.JSON(code, obj)
}
