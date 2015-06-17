package happyngine

import (
	"regexp"
)

type FormValidatorHandler func(*Context, FormElementInterface)

func RegexpFormValidator(pattern string) FormValidatorHandler {

	return func(c *Context, e FormElementInterface) {

		r := regexp.MustCompile(pattern)

		if !r.MatchString(e.FormValue()) {

			c.AddError(400, e.Error())
		}
	}
}
