package tools

import (
	validator "github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator(validate *validator.Validate) echo.Validator {
	return &CustomValidator{validator: validate}
}

func (v *CustomValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
