package validator

import (
    "errors"
)


type Validator struct {

    ValidatorHandler ValidatorHandler
    ErrorMessage string
}

type ValidatorHandler func(value string) bool

func New(validatorHandler ValidatorHandler, errorMessage string) *Validator {

    this := new(Validator)
    this.ValidatorHandler = validatorHandler
    this.ErrorMessage = errorMessage

    return this
}

func (this *Validator) IsValid(data string) error {

    if this.ValidatorHandler(data) {

        return nil
    }
    return errors.New(this.ErrorMessage)
}

func IsEqual(reference string) ValidatorHandler {

    return func(value string) bool {

        return reference == value
    }
}
