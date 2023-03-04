package http

import (
	"errors"
	"fmt"

	customValidation "github.com/Fajar-Islami/simple_manage_products/internal/validator"
	"github.com/go-playground/validator/v10"
)

// CustomValidator is struct used to create request validator
type CustomValidator struct {
	Validator    *validator.Validate
	ValidatorMap map[string]customValidation.ValidationHandler
}

type ValidatorParams interface {
	Validate() error
}

// NewValidator is function to create custom validator struct
func NewValidator() *CustomValidator {
	validate := validator.New()

	validatorMap := map[string]customValidation.ValidationHandler{}

	// Add new validation handlers here
	handlers := []customValidation.ValidationHandler{}

	for _, handler := range handlers {
		validatorMap[handler.Tag] = handler
		validate.RegisterValidation(handler.Tag, handler.Func)
	}

	return &CustomValidator{
		Validator:    validate,
		ValidatorMap: validatorMap,
	}
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.Validator.Struct(i); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range validationErrors {
				var message string
				if handler, ok := c.ValidatorMap[err.Tag()]; ok {
					message = fmt.Sprintf("validation error for '%s', Message: %s", err.Field(), handler.ErrorMessage)
				} else {
					message = fmt.Sprintf("validation error for '%s', Tag: %s", err.Field(), err.Tag())
				}
				if err.Param() != "" {
					message = fmt.Sprintf("%s %s", message, err.Param())
				}

				return errors.New(message)
			}
		}

		return err
	}

	if payload, ok := i.(ValidatorParams); ok {
		if err := payload.Validate(); err != nil {
			return err
		}
	}

	return nil
}
