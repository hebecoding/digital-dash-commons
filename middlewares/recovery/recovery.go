package recovery

import (
	"digital-dash-commons/rest"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Config struct {
	Filter func(c *fiber.Ctx) bool
	Logger *zap.Logger
}

func RecoverFromPanic(config ...Config) fiber.Handler {

	cfg := Config{}

	if len(config) > 0 {
		cfg = config[0]
	} else {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Print(err)
		}

		cfg.Logger = logger
	}

	return func(c *fiber.Ctx) error {

		defer func() {
			if err := recover(); err != nil {
				cfg.Logger.Debug("FATAL ERROR:", zap.Stack("message"))
				c.Status(http.StatusInternalServerError).JSON(rest.Response{
					Headers: c.GetRespHeaders(),
					Status:  http.StatusInternalServerError,
					Data:    c.Request().String(),
				})
			}
		}()

		return c.Next()
	}
}
