package usecase

import (
	"net/http"

	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/labstack/echo/v4"
)

func usecaseValidation(ctx echo.Context, params any) *helper.ErrorStruct {
	if err := ctx.Validate(params); err != nil {
		return &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}
	return nil
}
