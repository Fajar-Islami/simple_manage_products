package daos

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		FullName   string    `gorm:"not null;index"`
		Username   string    `gorm:"unique;not null"`
		Password   string    `gorm:"not null"`
		FirstOrder time.Time `gorm:"not null"`
	}

	FilterBooks struct {
		Limit, Offset int
		Fullname      string
	}
)
