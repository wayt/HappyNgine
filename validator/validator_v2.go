package validator

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"log"
	// "regexp"
)

func init() {
	// govalidator.ParamTagRegexMap["arraylength"] = regexp.MustCompile("^arraylength\\((\\d+)\\|(\\d+)\\)$")
	// govalidator.TagMap["arraylength"] = govalidator.Validator(func(a []interface{}, min, max int) bool {
	// 	l := len(a)
	//
	// 	if l >= min && l <= max {
	// 		return true
	// 	}
	//
	// 	return false
	// })
}

type Error struct {
	Label string `json:"label"`
	Field string `json:"field"`
	Text  string `json:"text"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func (e Errors) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func Validate(v interface{}) error {
	if ok, err := govalidator.ValidateStruct(v); !ok {
		m := FormatErrors(err)
		log.Println(m)
		return m
	}
	return nil
}

func FormatErrors(errs error) *Errors {

	e := new(Errors)
	e.Errors = make([]Error, 0)

	m := govalidator.ErrorsByField(errs)
	for key, value := range m {

		key = govalidator.CamelCaseToUnderscore(key)
		e.Errors = append(e.Errors, Error{
			Label: fmt.Sprintf("invalid_%s", key),
			Field: key,
			Text:  value,
		})
	}

	return e
}
