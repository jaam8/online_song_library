package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggingMiddleware добавляет логирование для каждого запроса
func LoggingMiddleware(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			log.Info("Request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("query", req.URL.RawQuery),
			)
			err := next(c)

			log.Info("Response",
				zap.String("method", req.Method),
				zap.Int("status", c.Response().Status),
				zap.String("path", req.URL.Path),
			)

			return err
		}
	}
}
