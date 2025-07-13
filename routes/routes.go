package routes

import (
	"JWT-Authentication-go/controllers"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/hello", controllers.Hello)
}
