package main

import (
	"backend/accounts"
	"backend/db"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	app := fiber.New()

	app.Use("/", static.New("./public"))

	db.InitDB()

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("PONG")
	})

	accounts.Endpoints(app)

	app.Listen(":8000")
}
