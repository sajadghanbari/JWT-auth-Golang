package main

import (
	"JWT-Authentication-go/api/routes"
	"JWT-Authentication-go/config"
	db "JWT-Authentication-go/data/database"
	_ "JWT-Authentication-go/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func main() {
	var cfg = config.GetConfig()
	db.InitDb(cfg)
	defer db.CloseDb()
	app := fiber.New()

	app.Use(logger.New())
	routes.InitRoutes(app)
	app.Get("/swagger/*", swagger.HandlerDefault)
	// پورت را از فایل کانفیگ بخوانید تا مدیریت آن راحت‌تر باشد
	app.Listen(":" + cfg.Server.Port)
}
