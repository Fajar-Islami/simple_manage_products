package usecase

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/dtos"
	"github.com/Fajar-Islami/simple_manage_products/internal/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type AuthUseCase interface {
	LoginUser(ctx echo.Context, params dtos.LoginRequest) (res dtos.LoginResp, err *helper.ErrorStruct)
	RegisterUser(ctx echo.Context, params dtos.RegisterRequest) (res string, err *helper.ErrorStruct)
}

type authUseCaseImpl struct {
	authrepository daos.AuthRepository
}

func NewAuthUseCase(authrepository daos.AuthRepository) AuthUseCase {
	return &authUseCaseImpl{
		authrepository: authrepository,
	}

}

func (oriu *authUseCaseImpl) LoginUser(ctx echo.Context, params dtos.LoginRequest) (res dtos.LoginResp, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()
	log := ctx.Logger()
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}

	// Check data from mysql
	resRepo, errRepo := oriu.authrepository.LoginUser(contx, params.Username)
	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		log.Error(errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Unauthorized"),
		}
	}
	if errRepo != nil {
		log.Error(errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusUnauthorized,
			Err:  errRepo,
		}
	}

	isValid := utils.CheckPasswordHash(params.Password, resRepo.Password)
	if !isValid {
		return res, &helper.ErrorStruct{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Unauthorized"),
		}
	}

	token, errGenerateToken := oriu.generateToken(resRepo)
	if errGenerateToken != nil {
		log.Error(errGenerateToken)
		return res, &helper.ErrorStruct{
			Code: http.StatusUnauthorized,
			Err:  errGenerateToken,
		}
	}

	res = dtos.LoginResp{
		Fullname: resRepo.FullName,
		Username: resRepo.Username,
		Token:    token,
	}

	return res, nil
}

func (oriu *authUseCaseImpl) RegisterUser(ctx echo.Context, params dtos.RegisterRequest) (res string, err *helper.ErrorStruct) {
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}

	hashPass, errHash := utils.HashPassword(params.Password)
	if errHash != nil {
		log.Println(errHash)
		return "", &helper.ErrorStruct{
			Code: http.StatusInternalServerError,
			Err:  errHash,
		}
	}

	_, errRepo := oriu.authrepository.RegisterUser(ctx.Request().Context(), daos.User{
		FullName: params.Fullname,
		Username: params.Username,
		Password: hashPass,
	})

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return "succeed register user", nil
}

func (oriu *authUseCaseImpl) generateToken(params daos.User) (res string, err error) {
	claims := jwt.MapClaims{}
	claims["username"] = params.Username
	claims["id"] = params.ID
	claims["exp"] = time.Now().Add(48 * time.Hour).Unix()

	token, errGenerateToken := utils.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		return res, err
	}

	return token, nil
}
