package daos

import (
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
	}
)
