package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/dtos"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	GetAllUserProfile(ctx echo.Context) error
	GetUserByID(ctx echo.Context) error
	GetMyUserByID(ctx echo.Context) error
	UpdateUserProfileByID(ctx echo.Context) error
	DeleteUserProfileByID(ctx echo.Context) error
}

type userControllerImpl struct {
	usersusecase usecase.UserUseCase
}

func NewUserController(usersusecase usecase.UserUseCase) UserController {
	return &userControllerImpl{
		usersusecase: usersusecase,
	}
}

func (uco *userControllerImpl) GetAllUserProfile(ctx echo.Context) error {
	filter := new(dtos.FilterUsers)
	if err := ctx.Bind(filter); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := uco.usersusecase.GetAllUserProfile(ctx, *filter)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

func (uco *userControllerImpl) GetUserByID(ctx echo.Context) error {
	userid := ctx.Param("userid")
	orderItemsIdInt, errConv := strconv.Atoi(userid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	res, err := uco.usersusecase.GetUserByID(ctx, orderItemsIdInt)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

func (uco *userControllerImpl) GetMyUserByID(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}
	res, err := uco.usersusecase.GetUserByID(ctx, userid)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

func (uco *userControllerImpl) UpdateUserProfileByID(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}
	params := new(dtos.UpdateUser)
	if err := ctx.Bind(params); err != nil {
		ctx.Logger().Error(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDUPDATEDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}

	res, err := uco.usersusecase.UpdateUserProfileByID(ctx, userid, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDUPDATEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDUPDATEDATA, "", res, http.StatusOK)
}

func (uco *userControllerImpl) DeleteUserProfileByID(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}

	res, err := uco.usersusecase.DeleteUserProfileByID(ctx, userid)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDDELETEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDDELETEDATA, "", res, http.StatusOK)
}
