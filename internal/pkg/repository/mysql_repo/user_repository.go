package mysql_repo

import (
	"context"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) daos.UsersRepository {
	return &UsersRepositoryImpl{
		db: db,
	}
}
func (ur *UsersRepositoryImpl) GetAllUserProfile(ctx context.Context, params daos.FilterUser) (res []daos.User, err error) {
	db := ur.db

	if err := db.Debug().WithContext(ctx).Limit(params.Limit).Offset(params.Offset).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (ur *UsersRepositoryImpl) GetMyUserByID(ctx context.Context, userid int) (res daos.User, err error) {
	if err := ur.db.First(&res, userid).WithContext(ctx).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (ur *UsersRepositoryImpl) UpdateUserProfileByID(ctx context.Context, userid int, data daos.User) (res string, err error) {
	var dataUsers daos.User
	if err = ur.db.Where("id = ? ", userid).First(&dataUsers).WithContext(ctx).Error; err != nil {
		return "Update user failed", gorm.ErrRecordNotFound
	}

	if err := ur.db.Model(dataUsers).Updates(&data).Where("id = ? ", userid).Error; err != nil {
		return "Update user failed", err
	}

	return res, nil
}

func (ur *UsersRepositoryImpl) DeleteUserProfileByID(ctx context.Context, userid int) (res string, err error) {
	var dataUsers daos.User
	if err = ur.db.Where("id = ?", userid).First(&dataUsers).WithContext(ctx).Error; err != nil {
		return "Delete user failed", gorm.ErrRecordNotFound
	}

	if err := ur.db.Model(dataUsers).Delete(&dataUsers).Error; err != nil {
		return "Delete user failed", err
	}

	return res, nil
}
