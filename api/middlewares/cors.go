package middlewares

import (
	"JWT-Authentication-go/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(cfg *config.Config) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5005", 
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,UPDATE",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
		MaxAge:           21600, 
	})
}
