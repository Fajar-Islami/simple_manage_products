package handler

import (
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"github.com/labstack/echo/v4"

	orderitemscontroller "github.com/Fajar-Islami/simple_manage_products/internal/pkg/controller"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/repository/mysql_repo"
	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/usecase"
)

func OrderItemsRoute(r *echo.Group, containerConf *container.Container) {
	repo := mysql_repo.NewOrderItemsRepository(containerConf.Mysqldb)
	usecase := usecase.NewOrderItemsUseCase(repo)
	controller := orderitemscontroller.NewOrderItemsController(usecase)

	orderItemsAPI := r.Group("/orderitems")
	orderItemsAPI.GET("", controller.GetAllOrderItems)
	orderItemsAPI.GET("/:orderitemsid", controller.GetOrderItemsByID)
	orderItemsAPI.POST("", controller.CreateOrderItems)
	orderItemsAPI.PUT("/:orderitemsid", controller.UpdateOrderItemsByID)
	orderItemsAPI.DELETE("/:orderitemsid", controller.DeleteOrderItemsByID)
}
