package http

import (
	"net/http"
	"strings"

	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/utils"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearer := c.Request().Header.Get("token")
		if bearer == "" {
			return helper.BuildResponse(c, false, "UNATHORIZED", "FAILED TO GET TOKEN", nil, http.StatusUnauthorized)
		}

		s := strings.Split(bearer, " ")
		if len(s) < 2 {
			return helper.BuildResponse(c, false, "UNATHORIZED", "FAILED TO GET TOKEN", nil, http.StatusUnauthorized)
		}

		claims, err := utils.DecodeToken(s[1])
		if err != nil {
			return helper.BuildResponse(c, false, "UNATHORIZED", err.Error(), nil, http.StatusUnauthorized)
		}

		c.Set("userid", int(claims.ID))

		// Go to next middleware:
		return next(c)
	}
}
