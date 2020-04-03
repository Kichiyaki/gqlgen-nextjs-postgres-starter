package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			latency := stop.Sub(start).String()
			logrus.WithFields(logrus.Fields{
				"remote_ip":  c.RealIP(),
				"host":       req.Host,
				"method":     req.Method,
				"uri":        req.RequestURI,
				"status":     res.Status,
				"user-agent": req.UserAgent(),
				"latency":    latency,
			}).Info("New request")
			return nil
		}
	}
}
