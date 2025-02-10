package validator

import (
	"github.com/go-playground/validator/v10"
)

var (
	defaultValidator *validator.Validate
)

func init()  {
	defaultValidator = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(i interface{}) error {
	return defaultValidator.Struct(i)
}
