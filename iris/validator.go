package iris

import (
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

func initValidator() {
	validate = validator.New()
}

func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}
