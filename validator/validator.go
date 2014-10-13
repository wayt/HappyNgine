package validator

import (
    "errors"
    "regexp"
    "time"
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

func IsEqual(references ...string) ValidatorHandler {

    return func(value string) bool {

        for _, reference := range references {

            if reference == value {

                return true
            }
        }
        return false
    }
}

func IsNotEmpty() ValidatorHandler {

    return func(value string) bool {

        if len(value) > 0 {

            return true
        }
        return false
    }
}

func Regexp(pattern string) ValidatorHandler {

    return func(value string) bool {

        r, _ := regexp.Compile(pattern)

        return r.MatchString(value)
    }
}

func IsEmail() ValidatorHandler {

    return Regexp(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]+`)
}

func IsDate() ValidatorHandler {

    return func(value string) bool {

        _, err := time.Parse("2006-01-02", value)
        if err != nil {

            return false
        }
        return true
    }
}
