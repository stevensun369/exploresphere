package accounts

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func postLogin(c fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.Unmarshal(c.Body(), &body)

	_, account := models.CheckAccount(body.Email)

	if bcrypt.CompareHashAndPassword(
		[]byte(account.Password), []byte(body.Password)) != nil {
		return utils.Error(c, errors.New("Codul introdus nu este corect"))
	}

	token := account.GenToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}
