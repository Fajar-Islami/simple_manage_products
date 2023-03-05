package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"github.com/Fajar-Islami/simple_manage_products/internal/server/http/handler"
	"github.com/Fajar-Islami/simple_manage_products/internal/utils"
)

func HTTPRouteInit(cont *container.Container, containerConf *container.Container) {
	e := echo.New()
	e.Validator = NewValidator()

	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(LoggerMiddleware(&containerConf.Logger.Log))

	utils.SecretKey = containerConf.Apps.SecretJwt

	api := e.Group("/api/v1") // /api
	handler.OrderItemsRoute(api, containerConf)
	handler.AuthRoute(api, containerConf)
	handler.UserRoute(api, containerConf, AuthMiddleware)
	handler.OrderHistoryRoute(api, containerConf, AuthMiddleware)

	port := fmt.Sprintf("%s:%d", containerConf.Apps.Host, containerConf.Apps.HttpPort)
	e.Logger.Fatal(e.Start(port))

}
