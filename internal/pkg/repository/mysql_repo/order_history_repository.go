package mysql_repo

import (
	"context"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"gorm.io/gorm"
)

type OrderHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderHistoryRepository(db *gorm.DB) daos.OrderHistoryRepository {
	return &OrderHistoryRepositoryImpl{
		db: db,
	}
}
func (alr *OrderHistoryRepositoryImpl) GetAllOrderHistory(ctx context.Context, params daos.FilterOrderHistory, userid int) (res []daos.OrderHistory, err error) {
	db := alr.db

	if err := db.Debug().WithContext(ctx).Limit(params.Limit).Offset(params.Offset).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (alr *OrderHistoryRepositoryImpl) GetOrderHistoryByID(ctx context.Context, historyid, userid int) (res daos.OrderHistory, err error) {
	if err := alr.db.First(&res, "id = ? and userid = ?", historyid, userid).WithContext(ctx).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (alr *OrderHistoryRepositoryImpl) CreateOrderHistory(ctx context.Context, data daos.OrderHistory, userid int) (res uint, err error) {
	// @@TODO use goroutine
	result := alr.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (alr *OrderHistoryRepositoryImpl) UpdateOrderHistoryByID(ctx context.Context, historyid, userid int, data daos.OrderHistory) (res string, err error) {
	var dataOrderHistory daos.OrderHistory
	if err = alr.db.Where("id = ? and userid = ?", historyid, userid).First(&dataOrderHistory).WithContext(ctx).Error; err != nil {
		return "Update order items failed", gorm.ErrRecordNotFound
	}

	if err := alr.db.Model(dataOrderHistory).Updates(&data).Where("id = ? and userid = ?", historyid, userid).Error; err != nil {
		return "Update order items failed", err
	}

	return res, nil
}

func (alr *OrderHistoryRepositoryImpl) DeleteOrderHistoryByID(ctx context.Context, historyid, userid string) (res string, err error) {
	var dataOrderHistory daos.OrderHistory
	if err = alr.db.Where("id = ? and userid = ?", historyid, userid).First(&dataOrderHistory).WithContext(ctx).Error; err != nil {
		return "Delete order items failed", gorm.ErrRecordNotFound
	}

	if err := alr.db.Model(dataOrderHistory).Delete(&dataOrderHistory).Error; err != nil {
		return "Delete order items failed", err
	}

	return res, nil
}
