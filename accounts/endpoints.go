package accounts

import (
	"backend/models"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

func Endpoints(app *fiber.App) {
	acc := app.Group("accounts")

	acc.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("PONG")
	})

	acc.Get("/whoami", func(c fiber.Ctx) error {
		var account models.Account
		utils.GetLocals(c, "account", &account)

		return c.JSON(account)
	}, models.AccountMiddleware)

	acc.Post("/login", postLogin)
	acc.Post("/signup", postSignup)

	acc.Post("/profile", postProfile, models.AccountMiddleware)
	acc.Post("/profilepic", postProfilePic, models.AccountMiddleware)

}
