package handler

import (
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"github.com/labstack/echo/v4"

	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/controller"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/repository/mysql_repo"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/usecase"
)

func AuthRoute(r *echo.Group, containerConf *container.Container) {
	repo := mysql_repo.NewAuthRepository(containerConf.Mysqldb)
	usecase := usecase.NewAuthUseCase(repo)
	controller := controller.NewAuthController(usecase)

	orderItemsAPI := r.Group("/auth")
	orderItemsAPI.POST("/login", controller.LoginUser)
	orderItemsAPI.POST("/register", controller.RegisterUser)
}
