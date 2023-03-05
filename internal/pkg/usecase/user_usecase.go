package usecase

import (
	"errors"
	"net/http"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/dtos"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type UserUseCase interface {
	GetAllUserProfile(ctx echo.Context, params dtos.FilterUsers) (res dtos.ResGetUser, err *helper.ErrorStruct)
	GetUserByID(ctx echo.Context, userid int) (res dtos.ResGetUserData, err *helper.ErrorStruct)
	UpdateUserProfileByID(ctx echo.Context, userid int, data dtos.UpdateUser) (res string, err *helper.ErrorStruct)
	DeleteUserProfileByID(ctx echo.Context, userid int) (res string, err *helper.ErrorStruct)
}

type userUseCaseImpl struct {
	userrepository daos.UsersRepository
}

func NewUserUseCase(userrepository daos.UsersRepository) UserUseCase {
	return &userUseCaseImpl{
		userrepository: userrepository,
	}

}

func (uc *userUseCaseImpl) GetAllUserProfile(ctx echo.Context, params dtos.FilterUsers) (res dtos.ResGetUser, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()
	log := ctx.Logger()
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}

	cpPage := params.Page
	dataRows := make([]dtos.ResGetUserData, 0)

	if params.Limit < 1 {
		params.Limit = 10
	}

	if params.Page < 1 {
		params.Page = 0
	} else {
		params.Page = (params.Page - 1) * params.Limit
	}

	resRepo, count, errRepo := uc.userrepository.GetAllUserProfile(contx, daos.FilterUser{
		Limit:    params.Limit,
		Offset:   params.Page,
		Fullname: params.Fullname,
	})
	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		log.Error(errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusUnauthorized,
			Err:  errors.New("No Data User"),
		}
	}
	if errRepo != nil {
		log.Error(errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusUnauthorized,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		dataRows = append(dataRows, dtos.ResGetUserData{
			Fullname: v.FullName,
			Username: v.Username,
		})
	}

	rows := params.Limit
	if rows > int(count) {
		rows = int(count)
	}

	res.Data = dataRows
	res.Page = cpPage
	res.Rows = rows
	res.TotalRows = int(count)
	return res, nil
}

func (uc *userUseCaseImpl) GetUserByID(ctx echo.Context, userid int) (res dtos.ResGetUserData, err *helper.ErrorStruct) {

	resRepo, errRepo := uc.userrepository.GetMyUserByID(ctx.Request().Context(), userid)

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New("user invalid"),
		}
	}

	return dtos.ResGetUserData{
		Fullname: resRepo.FullName,
		Username: resRepo.Username,
	}, nil
}

func (uc *userUseCaseImpl) UpdateUserProfileByID(ctx echo.Context, userid int, params dtos.UpdateUser) (res string, err *helper.ErrorStruct) {
	resRepo, errRepo := uc.userrepository.UpdateUserProfileByID(ctx.Request().Context(), userid, daos.User{
		FullName: params.Fullname,
		Username: params.Username,
		Password: params.Password,
	})

	if helper.MysqlCheckErrDuplicateEntry(errRepo) {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New("username already taken"),
		}
	}
	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New("user invalid"),
		}
	}

	return resRepo, nil
}

func (uc *userUseCaseImpl) DeleteUserProfileByID(ctx echo.Context, userid int) (res string, err *helper.ErrorStruct) {
	resRepo, errRepo := uc.userrepository.DeleteUserProfileByID(ctx.Request().Context(), userid)
	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New("user invalid"),
		}
	}

	return resRepo, nil
}
