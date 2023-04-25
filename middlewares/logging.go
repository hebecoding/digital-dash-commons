package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"log"
	"time"
)

const (
	reqID = "x-request-id"
)

func HTTPLogger(opts ...Option) fiber.Handler {
	cfg := &Config{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Logger == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Print(err)
		}
		cfg.Logger = logger
	}

	return func(c *fiber.Ctx) error {
		if len(cfg.Filters) > 0 {
			skip := true
			for _, filter := range cfg.Filters {
				if !filter(c) {
					skip = false
					break
				}
			}

			if skip {
				return c.Next()
			}
		}

		start := time.Now()
		requestID := c.Get(reqID)

		if requestID == "" {
			requestID = xid.New().String()
			c.Set(reqID, requestID)
		}

		err := c.Next()

		latency := time.Since(start).Milliseconds()
		statusCode := c.Response().StatusCode()

		logFields := []zap.Field{
			zap.Int("statusCode", statusCode),
			zap.Int64("latency", latency),
			zap.String("method", c.Method()),
			zap.String("uri", c.Path()),
			zap.String("ipAddress", c.IP()),
			zap.String("requestID", requestID),
		}

		switch {
		case statusCode >= 500:
			cfg.Logger.Error("Internal Server Error", logFields...)
		case statusCode >= 400:
			cfg.Logger.Warn("Client Error", logFields...)
		default:
			cfg.Logger.Info("Request Processed", logFields...)
		}

		return err
	}
}
