package happyngine

import (
	"regexp"
	"strconv"
	"time"
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

func IsEqual(refs ...string) FormValidatorHandler {

	return func(c *Context, e FormElementInterface) {

		for _, ref := range refs {
			if ref == e.FormValue() {
				return
			}
		}

		c.AddError(400, e.Error())
	}
}

func IsEmail() FormValidatorHandler {

	return RegexpFormValidator(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]+$`)
}

func IsInteger() FormValidatorHandler {

	return func(c *Context, e FormElementInterface) {

		RegexpFormValidator(`^(-|)[0-9]+$`)(c, e)
		if c.HasErrors() {
			return
		}

		// We don't check errors because of the previous regexp
		i, _ := strconv.ParseInt(e.FormValue(), 10, 64)
		e.SetValue(i)
	}
}

func IsUInteger() FormValidatorHandler {

	return func(c *Context, e FormElementInterface) {

		RegexpFormValidator(`^[0-9]+$`)(c, e)
		if c.HasErrors() {
			return
		}

		// We don't check errors because of the previous regexp
		i, _ := strconv.ParseUint(e.FormValue(), 10, 64)
		e.SetValue(i)
	}
}

func IsDate() FormValidatorHandler {

	return func(c *Context, e FormElementInterface) {

		_, err := time.Parse("2006-01-02", e.FormValue())
		if err != nil {
			c.AddError(400, e.Error())
		}
	}
}

func IsUUID() FormValidatorHandler {

	return RegexpFormValidator(`^[a-f\d]{8}(-[a-f\d]{4}){3}-[a-f\d]{12}?$`)
}
