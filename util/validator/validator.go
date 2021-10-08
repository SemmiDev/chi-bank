package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Struct(s interface{}) error {
	err := validator.New().Struct(s)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return NewErrFieldValidation(e)
		}
	}
	return nil
}

var ErrFieldValidation = errors.New("Field is not valid")

func NewErrFieldValidation(err validator.FieldError) error {
	return fmt.Errorf("%s: %w; format must be (%s=%s)", err.Field(), ErrFieldValidation, err.ActualTag(), err.Param())
}
