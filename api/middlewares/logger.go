package middlewares

import (
	"JWT-Authentication-go/config"
	logging "JWT-Authentication-go/pkg"

	"time"

	"github.com/gofiber/fiber/v2"
)

func DefaultStructuredLogger(cfg *config.Config) fiber.Handler {
	logger := logging.NewLogger(cfg)
	return structuredLogger(logger)
}

func structuredLogger(logger logging.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// --- Request Info ---
		path := c.OriginalURL()
		method := c.Method()
		clientIP := c.IP()
		bodyBytes := c.Body() // raw request body

		err := c.Next()
		// --- Response Info ---
		stop := time.Now()
		latency := stop.Sub(start)
		status := c.Response().StatusCode()
		respBody := c.Response().Body()

		keys := map[logging.ExtraKey]interface{}{}
		keys[logging.Path] = path
		keys[logging.ClientIp] = clientIP
		keys[logging.Method] = method
		keys[logging.Latency] = latency
		keys[logging.StatusCode] = status
		keys[logging.ErrorMessage] = ""
		keys[logging.BodySize] = len(respBody)
		keys[logging.RequestBody] = string(bodyBytes)
		keys[logging.ResponseBody] = string(respBody)

		logger.Info(logging.RequestResponse, logging.Api, "", keys)

		return err
	}
}
