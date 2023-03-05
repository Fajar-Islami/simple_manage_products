package controller

import (
	"log"
	"net/http"

	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/dtos"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type AuthController interface {
	LoginUser(ctx echo.Context) error
	RegisterUser(ctx echo.Context) error
}

type authControllerImpl struct {
	authsusecase usecase.AuthUseCase
}

func NewAuthController(authsusecase usecase.AuthUseCase) AuthController {
	return &authControllerImpl{
		authsusecase: authsusecase,
	}
}

func (aco *authControllerImpl) LoginUser(ctx echo.Context) error {
	params := new(dtos.LoginRequest)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := aco.authsusecase.LoginUser(ctx, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, "", res, http.StatusOK)
}

func (aco *authControllerImpl) RegisterUser(ctx echo.Context) error {
	params := new(dtos.RegisterRequest)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := aco.authsusecase.RegisterUser(ctx, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, "", res, http.StatusOK)
}
