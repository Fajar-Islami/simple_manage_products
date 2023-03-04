package http

import (
	"bytes"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func LoggerMiddleware(app *echo.Echo, log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			reqBody, err := readReqBody(c)
			if err != nil {
				return err
			}

			app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
				LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
					log.Info().
						Str("URI", v.URI).
						Int("status", v.Status).
						Str("user-agent", v.UserAgent).
						Interface("headers", v.Headers).
						Bytes("request body", reqBody).
						Msg("request")

					return nil
				},
			}))
			return
		}
	}
}

// func readResponseBody(c echo.Context) *bytes.Buffer {
// 	resBody := new(bytes.Buffer)
// 	mw := io.MultiWriter(c.Response().Writer, resBody)
// 	writer := &responseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
// 	c.Response().Writer = writer
// 	return resBody
// }

func readReqBody(c echo.Context) (body []byte, err error) {
	if c.Request().Body != nil { // Read
		body, err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return body, err
		}
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, err
}
