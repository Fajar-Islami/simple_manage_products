package daos

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		FullName   string `gorm:"not null;index"`
		Username   string `gorm:"unique;not null"`
		Password   string `gorm:"not null"`
		FirstOrder *time.Time
	}

	FilterUser struct {
		Limit, Offset int
		Fullname      string
	}
)

type UsersRepository interface {
	GetAllUserProfile(ctx context.Context, params FilterUser) (res []User, count int64, err error)
	GetMyUserByID(ctx context.Context, userid int) (res User, err error)
	UpdateUserProfileByID(ctx context.Context, userid int, data User) (res string, err error)
	DeleteUserProfileByID(ctx context.Context, userid int) (res string, err error)
}
