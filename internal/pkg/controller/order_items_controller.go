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

type OrderItemsController interface {
	GetAllOrderItems(ctx echo.Context) error
	GetOrderItemsByID(ctx echo.Context) error
	CreateOrderItems(ctx echo.Context) error
	UpdateOrderItemsByID(ctx echo.Context) error
	DeleteOrderItemsByID(ctx echo.Context) error
}

type orderItemsControllerImpl struct {
	orderitemsusecase usecase.OrderItemsUseCase
}

func NewOrderItemsController(orderitemsusecase usecase.OrderItemsUseCase) OrderItemsController {
	return &orderItemsControllerImpl{
		orderitemsusecase: orderitemsusecase,
	}
}

func (orico *orderItemsControllerImpl) GetAllOrderItems(ctx echo.Context) error {
	filter := new(dtos.FilterOrderItems)
	if err := ctx.Bind(filter); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := orico.orderitemsusecase.GetAllOrderItems(ctx, *filter)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

func (orico *orderItemsControllerImpl) GetOrderItemsByID(ctx echo.Context) error {
	orderitemsid := ctx.Param("orderitemsid")
	orderItemsIdInt, errConv := strconv.Atoi(orderitemsid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	res, err := orico.orderitemsusecase.GetOrderItemsByID(ctx, orderItemsIdInt)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}
func (orico *orderItemsControllerImpl) CreateOrderItems(ctx echo.Context) error {
	params := new(dtos.ReqDataOrderItems)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := orico.orderitemsusecase.CreateOrderItems(ctx, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, "", res, http.StatusOK)
}
func (orico *orderItemsControllerImpl) UpdateOrderItemsByID(ctx echo.Context) error {
	orderitemsid := ctx.Param("orderitemsid")
	orderItemsIdInt, errConv := strconv.Atoi(orderitemsid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDUPDATEDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	params := new(dtos.ReqDataUpdateOrderItems)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := orico.orderitemsusecase.UpdateOrderItemsByID(ctx, orderItemsIdInt, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDUPDATEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDUPDATEDATA, "", res, http.StatusOK)
}
func (orico *orderItemsControllerImpl) DeleteOrderItemsByID(ctx echo.Context) error {
	orderitemsid := ctx.Param("orderitemsid")
	orderItemsIdInt, errConv := strconv.Atoi(orderitemsid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDDELETEDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	res, err := orico.orderitemsusecase.DeleteOrderItemsByID(ctx, orderItemsIdInt)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDDELETEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDUPDATEDATA, "", res, http.StatusOK)
}
