package logging

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"log"
	"time"
)

type Config struct {
	Filter func(c *fiber.Ctx) bool
	Logger *zap.Logger
}

const (
	reqID = "X-Request-ID"
)

type loggedRequest struct {
	StatusCode      int               `json:"statusCode"`
	Latency         time.Duration     `json:"latency,omitempty"`
	Body            string            `json:"body,omitempty"`
	ResponseBody    string            `json:"responseBody,omitempty"`
	Method          string            `json:"method"`
	URI             string            `json:"uri"`
	IPAddress       string            `json:"ipAddress"`
	RequestID       string            `json:"requestID"`
	Query           string            `json:"queryString"`
	Headers         map[string]string `json:"headers"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
}

func HTTPLogger(config ...Config) fiber.Handler {

	start := time.Now()
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

		requestID := c.Get(reqID)

		if requestID == "" {
			requestID = xid.New().String()
			c.Set(reqID, requestID)
		}

		err := c.Next()

		req := loggedRequest{
			StatusCode:      c.Response().StatusCode(),
			URI:             c.Path(),
			Headers:         c.GetReqHeaders(),
			ResponseHeaders: c.GetRespHeaders(),
			Query:           c.Request().URI().QueryArgs().String(),
			Body:            string(c.Body()),
			ResponseBody:    c.Response().String(),
			IPAddress:       c.IP(),
			RequestID:       requestID,
			Method:          c.Method(),
			Latency:         time.Since(start),
		}

		cfg.Logger.Info(req.RequestID, zap.Any("request", req))

		return err
	}

}
