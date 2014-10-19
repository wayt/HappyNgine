package happy

import (
    "github.com/gohappy/happy/validator"
)

type Parameter struct {

    Name string
    Validators []*validator.Validator
    Required bool
}

func NewParameter(name string, required bool, validators ...*validator.Validator) *Parameter {

    this := new(Parameter)
    this.Name = name
    this.Validators = validators
    this.Required = required

    return this
}

func (this *Parameter) IsValid(data string) error {

    if len(data) == 0 && !this.Required {
        return nil
    }

    for _, validator := range this.Validators {

        err := validator.IsValid(data)
        if err != nil {

            return err
        }
    }
    return nil
}
