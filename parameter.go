package happy

import (
    "github.com/gohappy/happy/validator"
)

type Parameter struct {

    Name string
    Validators []*validator.Validator
}

func NewParameter(name string, validators ...*validator.Validator) *Parameter {

    this := new(Parameter)
    this.Name = name
    this.Validators = validators

    return this
}

func (this *Parameter) IsValid(data string) error {

    for _, validator := range this.Validators {

        err := validator.IsValid(data)
        if err != nil {

            return err
        }
    }
    return nil
}
