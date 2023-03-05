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

type OrderHistoryController interface {
	GetAllOrderHistory(ctx echo.Context) error
	GetOrderHistoryByID(ctx echo.Context) error
	CreateOrderHistory(ctx echo.Context) error
	UpdateOrderHistoryByID(ctx echo.Context) error
	DeleteOrderHistoryByID(ctx echo.Context) error
}

type orderHistoryControllerImpl struct {
	orderHistoryusecase usecase.OrderHistoryUseCase
}

func NewOrderHistoryController(orderHistoryusecase usecase.OrderHistoryUseCase) OrderHistoryController {
	return &orderHistoryControllerImpl{
		orderHistoryusecase: orderHistoryusecase,
	}
}

func (ohco *orderHistoryControllerImpl) GetAllOrderHistory(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}
	filter := new(dtos.FilterOrderHistory)
	if err := ctx.Bind(filter); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.orderHistoryusecase.GetAllOrderHistory(ctx, userid, *filter)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

func (ohco *orderHistoryControllerImpl) GetOrderHistoryByID(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}
	orderhistoryid := ctx.Param("orderhistoryid")
	orderHistoryIdInt, errConv := strconv.Atoi(orderhistoryid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.orderHistoryusecase.GetOrderHistoryByID(ctx, orderHistoryIdInt, userid)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}
func (ohco *orderHistoryControllerImpl) CreateOrderHistory(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}
	params := new(dtos.ReqCreateDataOrderHistoryItem)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.orderHistoryusecase.CreateOrderHistory(ctx, userid, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, "", res, http.StatusOK)
}
func (ohco *orderHistoryControllerImpl) UpdateOrderHistoryByID(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}
	orderhistoryid := ctx.Param("orderhistoryid")
	orderHistoryIdInt, errConv := strconv.Atoi(orderhistoryid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDUPDATEDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	params := new(dtos.ReqUpdateDataOrderHistoryItem)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.orderHistoryusecase.UpdateOrderHistoryByID(ctx, orderHistoryIdInt, userid, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDUPDATEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDUPDATEDATA, "", res, http.StatusOK)
}
func (ohco *orderHistoryControllerImpl) DeleteOrderHistoryByID(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		ctx.Logger().Error("TOKEN INVALID")
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, "TOKEN INVALID", nil, http.StatusUnauthorized)
	}
	orderhistoryid := ctx.Param("orderhistoryid")
	orderHistoryIdInt, errConv := strconv.Atoi(orderhistoryid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDDELETEDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.orderHistoryusecase.DeleteOrderHistoryByID(ctx, orderHistoryIdInt, userid)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDDELETEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDUPDATEDATA, "", res, http.StatusOK)
}
