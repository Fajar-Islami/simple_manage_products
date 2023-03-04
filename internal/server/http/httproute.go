package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
)

func HTTPRouteInit(cont *container.Container, containerConf *container.Container) {
	e := echo.New()
	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(LoggerMiddleware(e, *cont.Logger))

	port := fmt.Sprintf("%s:%d", containerConf.Apps.Host, containerConf.Apps.HttpPort)
	e.Logger.Fatal(e.Start(port))

	// api := e.Group("/api/v1") // /api

	// route.BooksRoute(api, containerConf)
}
