package daos

import (
	"context"

	"gorm.io/gorm"
)

type (
	OrderHistory struct {
		gorm.Model
		Descriptions string     `gorm:"not null"`
		UserID       uint       `gorm:"not null"`
		User         User       `gorm:"foreignKey:UserID"`
		OrderItemID  uint       `gorm:"not null"`
		OrderItem    OrderItems `gorm:"foreignKey:OrderItemID"`
	}

	FilterOrderHistory struct {
		Limit, Offset int
		Descriptions  string
	}
)

type OrderHistoryRepository interface {
	GetAllOrderHistory(ctx context.Context, params FilterOrderHistory, userid int) (res []OrderHistory, err error)
	GetOrderHistoryByID(ctx context.Context, historyid, userid int) (res OrderHistory, err error)
	CreateOrderHistory(ctx context.Context, data OrderHistory, userid int) (res uint, err error)
	UpdateOrderHistoryByID(ctx context.Context, historyid, userid int, data OrderHistory) (res string, err error)
	DeleteOrderHistoryByID(ctx context.Context, historyid, userid string) (res string, err error)
}
