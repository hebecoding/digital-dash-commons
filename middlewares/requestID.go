package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
)

const reqIDKey = "x-request-id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(reqIDKey)

		// Generate a new request ID if not present in the header
		if requestID == "" {
			requestID = xid.New().String()
		}

		// Set the request ID in the header and store it in the context
		c.Set(reqIDKey, requestID)
		c.Locals(reqIDKey, requestID)

		return c.Next()
	}
}
