package mysql_repo

import (
	"context"
	"fmt"

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
func (ur *UsersRepositoryImpl) GetAllUserProfile(ctx context.Context, params daos.FilterUser) (res []daos.User, count int64, err error) {
	db := ur.db

	if params.Fullname != "" {
		db = db.Where("full_name like ?", fmt.Sprint("%", params.Fullname, "%"))
	}

	if err := db.WithContext(ctx).Order("full_name asc").Limit(params.Limit).Offset(params.Offset).Find(&res).Error; err != nil {
		return res, count, err
	}

	if err := db.WithContext(ctx).Order("full_name asc").Model(&res).Count(&count).Error; err != nil {
		return res, count, err
	}

	return res, count, nil
}

func (ur *UsersRepositoryImpl) GetMyUserByID(ctx context.Context, userid int) (res daos.User, err error) {
	if err := ur.db.WithContext(ctx).First(&res, userid).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (ur *UsersRepositoryImpl) UpdateUserProfileByID(ctx context.Context, userid int, data daos.User) (res string, err error) {
	var dataUsers daos.User
	if err = ur.db.Where("id = ? ", userid).WithContext(ctx).First(&dataUsers).Error; err != nil {
		return "Update user failed", gorm.ErrRecordNotFound
	}

	if err := ur.db.WithContext(ctx).Model(dataUsers).Updates(&data).Where("id = ? ", userid).Error; err != nil {
		return "Update user failed", err
	}

	return "update user suucceed", nil
}

func (ur *UsersRepositoryImpl) DeleteUserProfileByID(ctx context.Context, userid int) (res string, err error) {
	var dataUsers daos.User
	if err = ur.db.Where("id = ?", userid).WithContext(ctx).First(&dataUsers).Error; err != nil {
		return "Delete user failed", gorm.ErrRecordNotFound
	}

	if err := ur.db.WithContext(ctx).Model(dataUsers).Delete(&dataUsers).Error; err != nil {
		return "Delete user failed", err
	}

	return "delete user suucceed", nil
}
