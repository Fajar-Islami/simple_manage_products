package mysql_repo

import (
	"context"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) daos.AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}
func (alr *AuthRepositoryImpl) LoginUser(ctx context.Context, params daos.User) (res daos.User, err error) {
	if err := alr.db.WithContext(ctx).First(&res, "username = ? and password = ?", params.Username, params.Password).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (alr *AuthRepositoryImpl) CreateUser(ctx context.Context, params daos.User) (res uint, err error) {
	result := alr.db.WithContext(ctx).Create(&params)
	if result.Error != nil {
		return res, result.Error
	}

	return params.ID, nil
}
