package routes

import (
	"JWT-Authentication-go/api/handlers"
	db "JWT-Authentication-go/data/database"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	database := db.GetDb()
	users := app.Group("/users")
	users.Post("/create", handlers.CreateUser(database))
	users.Get("/get-users", handlers.GetAllUsers(database))
	users.Delete("/delete/:id", handlers.DeleteUser(database))
}
