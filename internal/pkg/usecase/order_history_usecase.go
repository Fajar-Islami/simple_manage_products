package usecase

import (
	"errors"
	"log"
	"net/http"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/dtos"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/repository/redis_repo"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type OrderHistoryUseCase interface {
	GetAllOrderHistory(ctx echo.Context, userid int, params dtos.FilterOrderHistory) (res dtos.ResDataOrderHistory, err *helper.ErrorStruct)
	GetOrderHistoryByID(ctx echo.Context, orderHistoryid, userid int) (res dtos.ResDataOrderHistoryItem, err *helper.ErrorStruct)
	CreateOrderHistory(ctx echo.Context, userid int, params dtos.ReqCreateDataOrderHistoryItem) (res int, err *helper.ErrorStruct)
	UpdateOrderHistoryByID(ctx echo.Context, orderHistoryid, userid int, params dtos.ReqUpdateDataOrderHistoryItem) (res string, err *helper.ErrorStruct)
	DeleteOrderHistoryByID(ctx echo.Context, orderHistoryid, userid int) (res string, err *helper.ErrorStruct)
}

type orderHistoryUseCaseImpl struct {
	orderHistoryrepository   daos.OrderHistoryRepository
	redisorhistoryrepository redis_repo.RedisOrderHistoryRepository
}

func NewOrderHistoryUseCase(orderHistoryrepository daos.OrderHistoryRepository, redisorhistoryrepository redis_repo.RedisOrderHistoryRepository) OrderHistoryUseCase {
	return &orderHistoryUseCaseImpl{
		orderHistoryrepository:   orderHistoryrepository,
		redisorhistoryrepository: redisorhistoryrepository,
	}

}

func (ohc *orderHistoryUseCaseImpl) GetAllOrderHistory(ctx echo.Context, userid int, params dtos.FilterOrderHistory) (res dtos.ResDataOrderHistory, err *helper.ErrorStruct) {
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}
	cpPage := params.Page
	dataRows := make([]dtos.ResDataOrderHistoryItem, 0)

	if params.Limit < 1 {
		params.Limit = 10
	}

	if params.Page < 1 {
		params.Page = 0
	} else {
		params.Page = (params.Page - 1) * params.Limit
	}

	resRepo, count, errRepo := ohc.orderHistoryrepository.GetAllOrderHistory(ctx.Request().Context(), userid, daos.FilterOrderHistory{
		Limit:        params.Limit,
		Offset:       params.Page,
		Descriptions: params.Description,
	})

	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return res, &helper.ErrorStruct{
			Code: http.StatusNotFound,
			Err:  errors.New("No Data OrderHistory"),
		}
	}

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		dataRows = append(dataRows, dtos.ResDataOrderHistoryItem{
			DtosModel:   dtos.DtosModel{ID: v.ID, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt},
			Description: v.Descriptions,
			OrderItem: dtos.ResDataOrderItemsData{
				DtosModel: dtos.DtosModel{
					ID:        v.OrderItem.ID,
					CreatedAt: v.OrderItem.CreatedAt,
					UpdatedAt: v.OrderItem.UpdatedAt,
				},
				Name:      v.OrderItem.Name,
				Price:     v.OrderItem.Price,
				ExpiredAt: *v.OrderItem.ExpiredAt,
			},
		})
	}

	rows := params.Limit
	if rows > int(count) {
		rows = int(count)
	}

	res.Data = dataRows
	res.Page = cpPage
	res.Rows = rows
	res.TotalRows = int(count)
	return res, nil
}

func (ohc *orderHistoryUseCaseImpl) GetOrderHistoryByID(ctx echo.Context, orderHistoryid, userid int) (res dtos.ResDataOrderHistoryItem, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()

	// Check data from redis
	data, errRedis := ohc.redisorhistoryrepository.GetOrderHistoryCtx(contx, orderHistoryid, userid)
	if data != nil {
		return *data, nil
	}

	log.Println("errRedis", errRedis)

	// Check data from mysql
	resRepo, errRepo := ohc.orderHistoryrepository.GetOrderHistoryByID(contx, orderHistoryid, userid)
	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return res, &helper.ErrorStruct{
			Code: http.StatusNotFound,
			Err:  errors.New("Data OrderHistory not found"),
		}
	}

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	res = dtos.ResDataOrderHistoryItem{
		DtosModel:   dtos.DtosModel{ID: resRepo.ID, CreatedAt: resRepo.CreatedAt, UpdatedAt: resRepo.UpdatedAt},
		Description: resRepo.Descriptions,
		OrderItem: dtos.ResDataOrderItemsData{
			DtosModel: dtos.DtosModel{
				ID:        resRepo.OrderItem.ID,
				CreatedAt: resRepo.OrderItem.CreatedAt,
				UpdatedAt: resRepo.OrderItem.UpdatedAt,
			},
			Name:      resRepo.OrderItem.Name,
			Price:     resRepo.OrderItem.Price,
			ExpiredAt: *resRepo.OrderItem.ExpiredAt,
		},
	}

	// Set data to redis
	errRedis = ohc.redisorhistoryrepository.SetOrderHistoryCtx(contx, &res, userid)
	log.Println("errRedis set redis", errRedis)

	return res, nil
}

func (ohc *orderHistoryUseCaseImpl) CreateOrderHistory(ctx echo.Context, userid int, params dtos.ReqCreateDataOrderHistoryItem) (res int, err *helper.ErrorStruct) {
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}

	resRepo, errRepo := ohc.orderHistoryrepository.CreateOrderHistory(ctx.Request().Context(), daos.OrderHistory{
		Descriptions: params.Description,
		UserID:       uint(userid),
		OrderItemID:  uint(params.OrderItemID),
	})

	if helper.MysqlCheckErrDuplicateEntry(errRepo) {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New("order item is already taken"),
		}
	}

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return resRepo, nil
}

func (ohc *orderHistoryUseCaseImpl) UpdateOrderHistoryByID(ctx echo.Context, orderHistoryid, userid int, params dtos.ReqUpdateDataOrderHistoryItem) (res string, err *helper.ErrorStruct) {
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}

	resRepo, errRepo := ohc.orderHistoryrepository.UpdateOrderHistoryByID(ctx.Request().Context(), orderHistoryid, userid, daos.OrderHistory{
		Descriptions: params.Description,
		OrderItemID:  uint(params.OrderItemID),
	})

	if helper.MysqlCheckErrDuplicateEntry(errRepo) {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New("order item is already taken"),
		}
	}

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	// Delete key in redis
	errRedis := ohc.redisorhistoryrepository.DeleteOrderHistoryCtx(ctx.Request().Context(), orderHistoryid, userid)
	log.Println("delete redis err : ", errRedis)

	return resRepo, nil
}

func (ohc *orderHistoryUseCaseImpl) DeleteOrderHistoryByID(ctx echo.Context, orderHistoryid, userid int) (res string, err *helper.ErrorStruct) {
	resRepo, errRepo := ohc.orderHistoryrepository.DeleteOrderHistoryByID(ctx.Request().Context(), orderHistoryid, userid)
	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	// Delete key in redis
	errRedis := ohc.redisorhistoryrepository.DeleteOrderHistoryCtx(ctx.Request().Context(), orderHistoryid, userid)
	log.Println("delete redis err : ", errRedis)

	return resRepo, nil
}
