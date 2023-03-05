package daos

import (
	"context"

	"gorm.io/gorm"
)

type (
	OrderHistory struct {
		gorm.Model
		Descriptions string     `gorm:"not null"`
		UserID       uint       `gorm:"not null;uniqueIndex:idx_order_history_user_order_item"`
		User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
		OrderItemID  uint       `gorm:"not null;uniqueIndex:idx_order_history_user_order_item"`
		OrderItem    OrderItems `gorm:"foreignKey:OrderItemID;constraint:OnDelete:CASCADE"`
	}

	FilterOrderHistory struct {
		Limit, Offset int
		Descriptions  string
	}
)

type OrderHistoryRepository interface {
	GetAllOrderHistory(ctx context.Context, userid int, params FilterOrderHistory) (res []OrderHistory, count int64, err error)
	GetOrderHistoryByID(ctx context.Context, historyid, userid int) (res OrderHistory, err error)
	CreateOrderHistory(ctx context.Context, data OrderHistory) (res int, err error)
	UpdateOrderHistoryByID(ctx context.Context, historyid, userid int, data OrderHistory) (res string, err error)
	DeleteOrderHistoryByID(ctx context.Context, historyid, userid int) (res string, err error)
}
