package mysql_repo

import (
	"context"
	"fmt"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"gorm.io/gorm"
)

type OrderItemsRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderItemsRepository(db *gorm.DB) daos.OrderItemsRepository {
	return &OrderItemsRepositoryImpl{
		db: db,
	}
}
func (oir *OrderItemsRepositoryImpl) GetAllOrderItems(ctx context.Context, params daos.FilterOrderItems) (res []daos.OrderItems, count int64, err error) {
	db := oir.db

	if params.Name != "" {
		db = db.Where("name like ?", fmt.Sprint("%", params.Name, "%"))
	}

	if params.PriceLessThan > 0 {
		db = db.Where("price < ?", params.PriceLessThan)
	}

	if params.PriceMoreThan > 0 {
		db = db.Where("price > ?", params.PriceMoreThan)
	}

	if !params.WithExpired {
		db = db.Where("expired_at > CURDATE()")
	}

	if err := db.WithContext(ctx).Order("expired_at desc").Limit(params.Limit).Offset(params.Offset).Find(&res).Error; err != nil {
		return res, count, err
	}

	if err := db.WithContext(ctx).Order("expired_at desc").Model(&res).Count(&count).Error; err != nil {
		return res, count, err
	}
	return res, count, nil
}

func (oir *OrderItemsRepositoryImpl) GetOrderItemsByID(ctx context.Context, orderItemsid int) (res daos.OrderItems, err error) {
	if err := oir.db.WithContext(ctx).First(&res, orderItemsid).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (oir *OrderItemsRepositoryImpl) CreateOrderItems(ctx context.Context, data daos.OrderItems) (res uint, err error) {
	result := oir.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (oir *OrderItemsRepositoryImpl) UpdateOrderItemsByID(ctx context.Context, orderItemsid int, data daos.OrderItems) (res string, err error) {
	var dataOrderItems daos.OrderItems
	// Get first
	if err = oir.db.WithContext(ctx).Where("id = ? ", orderItemsid).First(&dataOrderItems).Error; err != nil {
		return "Update Order Item failed", gorm.ErrRecordNotFound
	}

	if err := oir.db.WithContext(ctx).Model(dataOrderItems).Updates(&data).Where("id = ? ", orderItemsid).Error; err != nil {
		return "Update Order Item failed", err
	}

	return "Update Order Items succeed", nil
}

func (oir *OrderItemsRepositoryImpl) DeleteOrderItemsByID(ctx context.Context, orderItemsid int) (res string, err error) {
	var dataOrderItems daos.OrderItems

	// Get First
	if err = oir.db.WithContext(ctx).Where("id = ?", orderItemsid).First(&dataOrderItems).Error; err != nil {
		return "Delete Order Items failed", gorm.ErrRecordNotFound
	}

	if err := oir.db.WithContext(ctx).Model(dataOrderItems).Delete(&dataOrderItems).Error; err != nil {
		return "Delete Order Items failed", err
	}

	return "Delete Order Item succeed", nil
}
