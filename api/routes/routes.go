package routes

import (
	"JWT-Authentication-go/api/handlers"
	"JWT-Authentication-go/api/middlewares"
	"JWT-Authentication-go/config"
	db "JWT-Authentication-go/data/database"

	"github.com/gofiber/fiber/v2"

)

func InitRoutes(app *fiber.App) {
	var cfg = config.GetConfig()
	app.Use(middlewares.Cors(cfg))
	
	app.Use(middlewares.DefaultStructuredLogger(cfg))

	database := db.GetDb()
	users := app.Group("/users",middlewares.LimitByRequest())
	
	users.Post("/create", handlers.CreateUser(database))
	users.Get("/get-users", handlers.GetAllUsers(database))
	users.Delete("/delete/:id", handlers.DeleteUser(database))
	users.Put("/update/:id", handlers.UpdateUser(database))
}
