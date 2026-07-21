package validator

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func Validate(data any) error {
	return validate.Struct(data)
}
