package middlewares

import (
	"digital-dash-commons/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Config struct {
	Filters []func(c *fiber.Ctx) bool
	Logger  *zap.Logger
}

type Option func(*Config)

func WithFilter(filter func(c *fiber.Ctx) bool) Option {
	return func(c *Config) {
		c.Filters = append(c.Filters, filter)
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func RecoverFromPanic(opts ...Option) fiber.Handler {
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

		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("Recovered from panic: %v", err)
				requestID := c.Get(reqID)

				cfg.Logger.Error(errMsg,
					zap.String("request_id", requestID),
					zap.String("method", c.Method()),
					zap.String("path", c.Path()))

				c.Status(http.StatusInternalServerError).JSON(models.Response{
					Headers:    c.GetRespHeaders(),
					StatusCode: http.StatusInternalServerError,
					Body:       []byte(errors.New(errMsg).Error()),
				})
			}
		}()

		return c.Next()
	}
}
