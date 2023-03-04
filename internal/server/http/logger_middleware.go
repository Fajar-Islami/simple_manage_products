package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// func LoggerMiddleware(app *echo.Echo, log zerolog.Logger) echo.MiddlewareFunc {
// 	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
// 		LogURI:        true,
// 		LogStatus:     true,
// 		LogLatency:    true,
// 		LogURIPath:    true,
// 		LogFormValues: middleware.DefaultCORSConfig.AllowHeaders,
// 		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
// 			log.Info().
// 				Str("URI", v.URI).
// 				Str("URIPATH", v.URIPath).
// 				Int("status", v.Status).
// 				Str("user-agent", v.UserAgent).
// 				Interface("headers", v.Headers).
// 				Interface("body", v.FormValues).
// 				Interface("body", v.FormValues).
// 				Msg("request")

// 			return nil
// 		},
// 	})

// }

func LoggerMiddleware(log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			reqBody, err := hookRequest(c)
			if err != nil {
				log.Error().Err(err)
			}

			var jsonCompact bytes.Buffer
			err = json.Compact(&jsonCompact, reqBody)
			if err != nil {
				log.Error().Err(err)
			}

			logger := log.With().
				Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Str("method", c.Request().Method).
				Str("uri", c.Request().RequestURI).
				// Bytes("body", reqBody).
				RawJSON("raw body", jsonCompact.Bytes()).
				Str("remote_ip", c.RealIP()).
				Logger()

			// Add logger to context
			c.Set("logger", &logger)

			// // Log pre-request message
			logger.Info().Msg("Pre-request")

			return next(c)

		}
	}
}

func hookRequest(c echo.Context) (body []byte, err error) {
	if c.Request().Body != nil { // Read
		body, err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return body, err
		}
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, err
}
