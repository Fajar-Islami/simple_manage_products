package validator

import v10 "github.com/go-playground/validator/v10"

type ValidationHandler struct {
	Tag          string
	Func         v10.Func
	ErrorMessage string
}
