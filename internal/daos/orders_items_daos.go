package daos

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	OrderItems struct {
		gorm.Model
		Name      string     `gorm:"not null"`
		Price     int        `gorm:"not null"`
		ExpiredAt *time.Time `gorm:"default:null"`
	}

	FilterOrderItems struct {
		Limit, Offset                int
		Name                         string
		PriceMoreThan, PriceLessThan int
		WithExpired                  bool
	}
)

type OrderItemsRepository interface {
	GetAllOrderItems(ctx context.Context, params FilterOrderItems) (res []OrderItems, count int64, err error)
	GetOrderItemsByID(ctx context.Context, orderItemsid int) (res OrderItems, err error)
	CreateOrderItems(ctx context.Context, data OrderItems) (res uint, err error)
	UpdateOrderItemsByID(ctx context.Context, orderItemsid int, data OrderItems) (res string, err error)
	DeleteOrderItemsByID(ctx context.Context, orderItemsid int) (res string, err error)
}
