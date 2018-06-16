package gin

import (
	"gopkg.in/go-playground/validator.v9"
)

var (
	Validate *validator.Validate
)

func initValidator() {
	Validate = validator.New()
}

func ValidateStruct(obj interface{}) error {
	return Validate.Struct(obj)
}
