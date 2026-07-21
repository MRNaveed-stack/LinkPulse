package validator

import "github.com/go-playground/validator/v10"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   e.Field(),
				Message: e.Error(),
			})
		}
		return validationErrors
	}
	return []ValidationError{
		{
			Field:   "unknown",
			Message: err.Error(),
		},
	}
}
