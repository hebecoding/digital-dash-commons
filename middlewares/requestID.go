package middlewares

import (
	"digital-dash-commons/utils"
	"github.com/gofiber/fiber/v2"
)

const reqIDKey = "x-request-id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(reqIDKey)

		// Generate a new request ID if not present in the header
		if requestID == "" {
			requestID = string(utils.NewXID())
		}

		// Set the request ID in the header and store it in the context
		c.Set(reqIDKey, requestID)
		c.Locals(reqIDKey, requestID)

		return c.Next()
	}
}
