package mysql_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"gorm.io/gorm"
)

type OrderHistoryRepositoryImpl struct {
	db     *gorm.DB
	logger container.Logger
}

func NewOrderHistoryRepository(db *gorm.DB, logger container.Logger) daos.OrderHistoryRepository {
	return &OrderHistoryRepositoryImpl{
		db:     db,
		logger: logger,
	}
}
func (ohr *OrderHistoryRepositoryImpl) GetAllOrderHistory(ctx context.Context, userid int, params daos.FilterOrderHistory) (res []daos.OrderHistory, count int64, err error) {
	db := ohr.db

	if params.Descriptions != "" {
		db = db.Where("descriptions like ?", fmt.Sprint("%", params.Descriptions, "%"))
	}

	if err := db.WithContext(ctx).Order("created_at desc").Limit(params.Limit).Offset(params.Offset).Preload("OrderItem", "deleted_at IS NULL").Find(&res, "user_id = ?", userid).Error; err != nil {
		return res, count, err
	}

	if err := db.WithContext(ctx).Model(&res).Where("user_id = ?", userid).Count(&count).Error; err != nil {
		return res, count, err
	}

	return res, count, nil
}

func (ohr *OrderHistoryRepositoryImpl) GetOrderHistoryByID(ctx context.Context, historyid, userid int) (res daos.OrderHistory, err error) {
	if err := ohr.db.WithContext(ctx).Preload("OrderItem", "deleted_at IS NULL").First(&res, "id = ? and user_id = ?", historyid, userid).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (ohr *OrderHistoryRepositoryImpl) CreateOrderHistory(ctx context.Context, data daos.OrderHistory) (res int, err error) {
	var userData daos.User
	var orderItems daos.OrderItems

	db := ohr.db.WithContext(ctx).Begin()

	// @@TODO use goroutine

	// Check if user first oder
	if err := db.Debug().Select("id,first_order").Find(&userData, "id = ?", data.UserID).Error; err != nil {
		msg := "Get data user failed for checking is first order"
		if err == nil && userData.ID == 0 {
			err = fmt.Errorf(msg)
		}
		ohr.logger.Log.Error().Err(err).Stack().Msg("Get data user failed for checking is first order")
		db.Rollback()
		return res, gorm.ErrRecordNotFound
	}

	if userData.FirstOrder == nil {
		now := time.Now()
		userData.FirstOrder = &now
		if err := db.Debug().Model(userData).Update("first_order", now).Where("id = ?", data.UserID).Error; err != nil {
			ohr.logger.Log.Error().Err(err).Stack().Msg("update data user first order failed")
			db.Rollback()
			return res, err
		}
	}

	// Check if item not expire
	if err := db.Debug().Select("id,expired_at").First(&orderItems, "id = ? and expired_at > CURDATE()", data.OrderItemID).Error; err != nil || orderItems.ID == 0 {
		ohr.logger.Log.Error().Stack().Err(err).Msg("Data not found")
		db.Rollback()
		return res, gorm.ErrRecordNotFound
	}

	// Insert Data
	result := db.Debug().Create(&data)
	if result.Error != nil {
		ohr.logger.Log.Error().Stack().Err(err).Msg("")
		db.Rollback()
		return res, result.Error
	}

	db.Commit()

	return int(data.ID), nil
}

func (ohr *OrderHistoryRepositoryImpl) UpdateOrderHistoryByID(ctx context.Context, historyid, userid int, data daos.OrderHistory) (res string, err error) {
	var dataOrderHistory daos.OrderHistory
	var orderItems daos.OrderItems

	db := ohr.db.WithContext(ctx).Begin()

	// Check data if exist
	if err = db.Where("id = ? and user_id = ?", historyid, userid).First(&dataOrderHistory).Error; err != nil {
		ohr.logger.Log.Error().Stack().Err(err).Msg("Data not found")
		db.Rollback()
		return res, gorm.ErrRecordNotFound
	}

	// Check if item not expire
	if err := db.Select("id,expired_at").First(&orderItems, "id = ? and expired_at > CURDATE()", data.OrderItemID).Error; err != nil || orderItems.ID == 0 {
		ohr.logger.Log.Error().Stack().Err(err).Msg("Data not found")
		db.Rollback()
		return res, gorm.ErrRecordNotFound
	}

	// Update Data
	if err := db.Debug().Model(dataOrderHistory).Updates(&data).Where("id = ? and user_id = ?", historyid, userid).Error; err != nil {
		ohr.logger.Log.Error().Stack().Err(err).Msg("Update data failed")
		db.Rollback()
		return helper.UpdateDataFailed, err
	}

	db.Commit()

	return helper.UpdateDataSucceed, nil
}

func (ohr *OrderHistoryRepositoryImpl) DeleteOrderHistoryByID(ctx context.Context, historyid, userid int) (res string, err error) {
	var dataOrderHistory daos.OrderHistory
	if err = ohr.db.WithContext(ctx).Where("id = ? and user_id = ?", historyid, userid).First(&dataOrderHistory).Error; err != nil {
		return helper.DeleteDataFailed, gorm.ErrRecordNotFound
	}

	if err := ohr.db.WithContext(ctx).Model(dataOrderHistory).Delete(&dataOrderHistory).Error; err != nil {
		return helper.DeleteDataFailed, err
	}

	return helper.DeleteDataSucceed, nil
}
