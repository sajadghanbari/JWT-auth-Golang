package main

import (
	"JWT-Authentication-go/api/routes"
	"JWT-Authentication-go/config"
	db "JWT-Authentication-go/data/database"
	_ "JWT-Authentication-go/docs"
	"JWT-Authentication-go/pkg/logging"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"
)

func main() {
	
	var cfg = config.GetConfig()
	logger := logging.NewLogger(cfg)
	err := db.InitDb(cfg)
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	defer db.CloseDb()
	app := fiber.New()


	routes.InitRoutes(app)
	setupSwagger(app)
	app.Listen(":" + cfg.Server.Port)
}

func setupSwagger(app *fiber.App) {
    app.Get("/swagger/*", swagger.HandlerDefault)
}
