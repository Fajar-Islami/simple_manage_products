package usecase

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/dtos"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/repository/redis_repo"
	"github.com/Fajar-Islami/simple_manage_products/internal/utils"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type OrderItemsUseCase interface {
	GetAllOrderItems(ctx echo.Context, params dtos.FilterOrderItems) (res dtos.ResDataOrderItems, err *helper.ErrorStruct)
	GetOrderItemsByID(ctx echo.Context, orderItemsid int) (res dtos.ResDataOrderItemsData, err *helper.ErrorStruct)
	CreateOrderItems(ctx echo.Context, params dtos.ReqDataOrderItems) (res uint, err *helper.ErrorStruct)
	UpdateOrderItemsByID(ctx echo.Context, orderItemsid int, params dtos.ReqDataUpdateOrderItems) (res string, err *helper.ErrorStruct)
	DeleteOrderItemsByID(ctx echo.Context, orderItemsid int) (res string, err *helper.ErrorStruct)
}

type orderItemsUseCaseImpl struct {
	orderitemsrepository      daos.OrderItemsRepository
	redisorderitemsrepository redis_repo.RedisOrderItemsRepository
}

func NewOrderItemsUseCase(orderitemsrepository daos.OrderItemsRepository, redisorderitemsrepository redis_repo.RedisOrderItemsRepository) OrderItemsUseCase {
	return &orderItemsUseCaseImpl{
		orderitemsrepository:      orderitemsrepository,
		redisorderitemsrepository: redisorderitemsrepository,
	}

}

func (oriu *orderItemsUseCaseImpl) GetAllOrderItems(ctx echo.Context, params dtos.FilterOrderItems) (res dtos.ResDataOrderItems, err *helper.ErrorStruct) {
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}
	cpPage := params.Page
	dataRows := make([]dtos.ResDataOrderItemsData, 0)

	if params.Limit < 1 {
		params.Limit = 10
	}

	if params.Page < 1 {
		params.Page = 0
	} else {
		params.Page = (params.Page - 1) * params.Limit
	}

	if params.PriceLessThan > 0 && params.PriceMoreThan > 0 {
		if params.PriceLessThan > params.PriceMoreThan {
			return res, &helper.ErrorStruct{
				Code: http.StatusBadRequest,
				Err:  errors.New("price_less_than must be smaller than price_more_than"),
			}
		}
	}

	resRepo, count, errRepo := oriu.orderitemsrepository.GetAllOrderItems(ctx.Request().Context(), daos.FilterOrderItems{
		Limit:         params.Limit,
		Offset:        params.Page,
		PriceMoreThan: params.PriceMoreThan,
		PriceLessThan: params.PriceLessThan,
		Name:          params.Name,
		WithExpired:   params.WithExpired,
	})

	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return res, &helper.ErrorStruct{
			Code: http.StatusNotFound,
			Err:  errors.New("No Data OrderItems"),
		}
	}

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		dataRows = append(dataRows, dtos.ResDataOrderItemsData{
			DtosModel: dtos.DtosModel{
				ID:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			},
			Name:      v.Name,
			Price:     v.Price,
			ExpiredAt: *v.ExpiredAt,
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

func (oriu *orderItemsUseCaseImpl) GetOrderItemsByID(ctx echo.Context, orderItemsid int) (res dtos.ResDataOrderItemsData, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()

	// Check data from redis
	data, errRedis := oriu.redisorderitemsrepository.GetOrderItemsCtx(contx, orderItemsid)
	if data != nil {
		return *data, nil
	}

	log.Println("errRedis", errRedis)

	// Check data from mysql
	resRepo, errRepo := oriu.orderitemsrepository.GetOrderItemsByID(contx, orderItemsid)
	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return res, &helper.ErrorStruct{
			Code: http.StatusNotFound,
			Err:  errors.New("No Data OrderItems"),
		}
	}

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	res = dtos.ResDataOrderItemsData{
		DtosModel: dtos.DtosModel{
			ID:        resRepo.ID,
			CreatedAt: resRepo.CreatedAt,
			UpdatedAt: resRepo.UpdatedAt,
		},
		Name:      resRepo.Name,
		Price:     resRepo.Price,
		ExpiredAt: *resRepo.ExpiredAt,
	}

	// Set data to redis
	errRedis = oriu.redisorderitemsrepository.SetOrderItemsCtx(contx, &res)
	log.Println("errRedis set redis", errRedis)

	return res, nil
}

func (oriu *orderItemsUseCaseImpl) CreateOrderItems(ctx echo.Context, params dtos.ReqDataOrderItems) (res uint, err *helper.ErrorStruct) {
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}

	expireDate, errTgl := utils.ShortDateFromString(params.ExpiredAt)
	if errTgl != nil {
		log.Println(errTgl)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errTgl,
		}
	}

	if expireDate.Before(time.Now()) {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New("Expired date is already passed"),
		}
	}

	resRepo, errRepo := oriu.orderitemsrepository.CreateOrderItems(ctx.Request().Context(), daos.OrderItems{
		Name:      params.Name,
		Price:     params.Price,
		ExpiredAt: &expireDate,
	})
	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return resRepo, nil
}
func (oriu *orderItemsUseCaseImpl) UpdateOrderItemsByID(ctx echo.Context, orderItemsid int, params dtos.ReqDataUpdateOrderItems) (res string, err *helper.ErrorStruct) {
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}

	expireDate, errTgl := utils.ShortDateFromString(params.ExpiredAt)
	if errTgl != nil {
		log.Println(errTgl)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errTgl,
		}
	}

	resRepo, errRepo := oriu.orderitemsrepository.UpdateOrderItemsByID(ctx.Request().Context(), orderItemsid, daos.OrderItems{
		Name:      params.Name,
		Price:     params.Price,
		ExpiredAt: &expireDate,
	})

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	// Delete key in redis
	errRedis := oriu.redisorderitemsrepository.DeleteOrderItemsCtx(ctx.Request().Context(), orderItemsid)
	log.Println("delete redis err : ", errRedis)

	return resRepo, nil
}
func (oriu *orderItemsUseCaseImpl) DeleteOrderItemsByID(ctx echo.Context, orderItemsid int) (res string, err *helper.ErrorStruct) {
	resRepo, errRepo := oriu.orderitemsrepository.DeleteOrderItemsByID(ctx.Request().Context(), orderItemsid)
	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	// Delete key in redis
	errRedis := oriu.redisorderitemsrepository.DeleteOrderItemsCtx(ctx.Request().Context(), orderItemsid)
	log.Println("delete redis err : ", errRedis)

	return resRepo, nil
}
