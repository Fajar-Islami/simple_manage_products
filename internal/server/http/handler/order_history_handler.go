package handler

import (
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"github.com/Fajar-Islami/simple_manage_products/internal/utils"
	"github.com/labstack/echo/v4"

	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/controller"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/repository/mysql_repo"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/repository/redis_repo"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/usecase"
)

func OrderHistoryRoute(r *echo.Group, containerConf *container.Container, authMiddleware utils.MiddlewareType) {
	redisClient := redis_repo.NewRedisRepoOrderHistory(containerConf.Redis, &containerConf.Logger.Log)
	repo := mysql_repo.NewOrderHistoryRepository(containerConf.Mysqldb, *containerConf.Logger)
	usecase := usecase.NewOrderHistoryUseCase(repo, redisClient)
	controller := controller.NewOrderHistoryController(usecase)

	orderItemsAPI := r.Group("/orderhistory")
	orderItemsAPI.Use(echo.MiddlewareFunc(authMiddleware))
	orderItemsAPI.GET("", controller.GetAllOrderHistory)
	orderItemsAPI.GET("/:orderhistoryid", controller.GetOrderHistoryByID)
	orderItemsAPI.POST("", controller.CreateOrderHistory)
	orderItemsAPI.PUT("/:orderhistoryid", controller.UpdateOrderHistoryByID)
	orderItemsAPI.DELETE("/:orderhistoryid", controller.DeleteOrderHistoryByID)
}
