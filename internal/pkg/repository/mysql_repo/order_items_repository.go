package mysql_repo

import (
	"context"

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
func (oir *OrderItemsRepositoryImpl) GetAllOrderItems(ctx context.Context, params daos.FilterOrderItems) (res []daos.OrderItems, err error) {
	db := oir.db

	if params.Name != "" {
		db = db.Where("name like ?%%", params.Name)
	}

	if params.PriceLessThan > 0 {
		db = db.Where("price < ?", params.PriceLessThan)
	}

	if params.PriceMoreThan > 0 {
		db = db.Where("price > ?", params.PriceMoreThan)
	}

	if err := db.Debug().WithContext(ctx).Limit(params.Limit).Offset(params.Offset).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (oir *OrderItemsRepositoryImpl) GetOrderItemsByID(ctx context.Context, orderItemsid int) (res daos.OrderItems, err error) {
	if err := oir.db.First(&res, orderItemsid).WithContext(ctx).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (oir *OrderItemsRepositoryImpl) CreateOrderItems(ctx context.Context, data daos.OrderItems) (res uint, err error) {
	result := oir.db.Create(&data).WithContext(ctx)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (oir *OrderItemsRepositoryImpl) UpdateOrderItemsByID(ctx context.Context, orderItemsid int, data daos.OrderItems) (res string, err error) {
	var dataOrderItems daos.OrderItems
	// Get first
	if err = oir.db.Where("id = ? ", orderItemsid).First(&dataOrderItems).WithContext(ctx).Error; err != nil {
		return "Update Order Item failed", gorm.ErrRecordNotFound
	}

	if err := oir.db.Model(dataOrderItems).Updates(&data).Where("id = ? ", orderItemsid).Error; err != nil {
		return "Update Order Item failed", err
	}

	return "Update Order Items succeed", nil
}

func (oir *OrderItemsRepositoryImpl) DeleteOrderItemsByID(ctx context.Context, orderItemsid int) (res string, err error) {
	var dataOrderItems daos.OrderItems

	// Get First
	if err = oir.db.Where("id = ?", orderItemsid).First(&dataOrderItems).WithContext(ctx).Error; err != nil {
		return "Delete Order Items failed", gorm.ErrRecordNotFound
	}

	if err := oir.db.Model(dataOrderItems).Delete(&dataOrderItems).Error; err != nil {
		return "Delete Order Items failed", err
	}

	return "Delete Order Item succeed", nil
}
