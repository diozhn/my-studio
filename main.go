package main

import (
	"my-studio/database"
	"my-studio/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to My Studio!")
	})

	routes.RegisterArtworkRoutes(app)

	app.Static("/uploads", "./uploads")
	app.Listen(":3000")
}
