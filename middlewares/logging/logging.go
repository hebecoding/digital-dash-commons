package logging

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

type Config struct {
	Filter func(c *fiber.Ctx) bool
	Logger *zap.Logger
}

const (
	reqID         = "X-Request-ID"
	interactionID = "X-Interaction-ID"
)

type loggedRequest struct {
	StatusCode      int               `json:"statusCode"`
	Body            string            `json:"body,omitempty"`
	ResponseBody    string            `json:"responseBody,omitempty"`
	Method          string            `json:"method"`
	URI             string            `json:"uri"`
	IPAddress       string            `json:"ipAddress"`
	RequestID       string            `json:"requestID"`
	InteractionID   string            `json:"interactionID"`
	Latency         string            `json:"latency,omitempty"`
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
		intID := c.Get(interactionID)

		if requestID == "" {
			requestID = xid.New().String()
			c.Set(reqID, requestID)
		}

		if intID == "" {
			intID = xid.New().String()
			c.Set(interactionID, intID)
		}

		err := c.Next()

		req := loggedRequest{
			StatusCode:      c.Response().StatusCode(),
			Headers:         c.GetReqHeaders(),
			ResponseHeaders: c.GetRespHeaders(),
			Body:            string(c.Body()),
			ResponseBody:    c.Response().String(),
			IPAddress:       c.IP(),
			RequestID:       requestID,
			InteractionID:   intID,
			Method:          c.Method(),
			Latency:         strconv.FormatInt(time.Since(start).Milliseconds(), 10),
		}

		cfg.Logger.Info(fmt.Sprintf("RequestID: %s", req.RequestID), zap.Any("Request", req))

		return err
	}

}