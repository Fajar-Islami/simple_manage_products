package daos

import (
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
