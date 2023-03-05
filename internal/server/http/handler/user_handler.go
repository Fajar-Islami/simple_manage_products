package handler

import (
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"github.com/Fajar-Islami/simple_manage_products/internal/utils"
	"github.com/labstack/echo/v4"

	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/controller"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/repository/mysql_repo"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/usecase"
)

func UserRoute(r *echo.Group, containerConf *container.Container, authMiddleware utils.MiddlewareType) {
	repo := mysql_repo.NewUsersRepository(containerConf.Mysqldb)
	usecase := usecase.NewUserUseCase(repo)
	controller := controller.NewUserController(usecase)

	orderItemsAPI := r.Group("/user")
	orderItemsAPI.GET("/list", controller.GetAllUserProfile)
	orderItemsAPI.GET("/:userid", controller.GetUserByID)
	orderItemsAPI.GET("/my", authMiddleware(controller.GetMyUserByID))
	orderItemsAPI.PUT("", authMiddleware(controller.UpdateUserProfileByID))
	orderItemsAPI.DELETE("", authMiddleware(controller.DeleteUserProfileByID))
}
