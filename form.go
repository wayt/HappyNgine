package happyngine

type FormElementInterface interface {
	Name() string
	Validate(*Context)
	SetFormValue(string)
	FormValue() string
	SetValue(interface{})
	Value() interface{}
	Error() string
	Required() bool
}

type FormElement struct {
	name        string
	formValue   string
	value       interface{}
	errorString string
	required    bool
	handlers    []FormValidatorHandler
}

func NewFormElement(name, errStr string) *FormElement {

	return &FormElement{
		name:        name,
		formValue:   "",
		value:       nil,
		errorString: errStr,
		required:    true,
	}
}

func (e *FormElement) Name() string {

	return e.name
}

func (e *FormElement) Validate(c *Context) {

	for _, h := range e.handlers {
		h(c, e)
		if c.HasErrors() {
			break
		}
	}
}

func (e *FormElement) SetFormValue(v string) {
	e.formValue = v
}

func (e *FormElement) FormValue() string {
	return e.formValue
}

func (e *FormElement) SetValue(i interface{}) {
	e.value = i
}

func (e *FormElement) Value() interface{} {
	return e.value
}

func (e *FormElement) Error() string {
	return e.errorString
}

func (e *FormElement) Required() bool {
	return e.required
}

func (e *FormElement) SetRequired(r bool) {
	e.required = r
}

func (e *FormElement) AddValidator(h FormValidatorHandler) *FormElement {

	e.handlers = append(e.handlers, h)

	return e
}

type Form struct {
	Context  *Context
	Elements map[string]FormElementInterface
}

func NewForm(c *Context, elems ...FormElementInterface) *Form {

	f := &Form{
		Context: c,
	}

	f.Elements = make(map[string]FormElementInterface)

	for _, e := range elems {
		f.Elements[e.Name()] = e
	}

	return f
}

func (f *Form) AddElement(e FormElementInterface) *Form {
	f.Elements[e.Name()] = e
	return f
}

func (f *Form) Elem(name string) FormElementInterface {
	return f.Elements[name]
}

func (f *Form) fillElements() bool {

	for _, e := range f.Elements {

		if value := f.Context.GetParam(e.Name()); len(value) > 0 {
			e.SetFormValue(value)
		} else if e.Required() {
			f.Context.AddError(400, e.Error())
		}
	}

	return !f.Context.HasErrors()
}

func (f *Form) IsValid() bool {

	if !f.fillElements() {
		return false
	}

	for _, e := range f.Elements {

		e.Validate(f.Context)
	}

	return !f.Context.HasErrors()
}
