package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"github.com/Fajar-Islami/simple_manage_products/internal/server/http/handler"
)

func HTTPRouteInit(cont *container.Container, containerConf *container.Container) {
	e := echo.New()

	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(LoggerMiddleware(*containerConf.Logger))
	e.Validator = NewValidator()

	api := e.Group("/api/v1") // /api
	handler.OrderItemsRoute(api, containerConf)

	port := fmt.Sprintf("%s:%d", containerConf.Apps.Host, containerConf.Apps.HttpPort)
	e.Logger.Fatal(e.Start(port))

}
