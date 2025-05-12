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

func postSignup(c fiber.Ctx) error {
	// unmarshaling body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	json.Unmarshal(c.Body(), &body)

	// checking if it exists
	exists, account := models.CheckAccount(body.Email)

	// error if exists
	if exists {
		return utils.Error(c, errors.New("An account with this email adress already exists"))
		// if it doesn't exist, we create it:
	} else {
		account.Email = body.Email
		account.Name = body.Name

		// hashing the password
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(body.Password), 10,
		)
		if err != nil {
			return utils.Error(c, err)
		}
		account.Password = string(hashedPassword)

		// creating the account
		err = account.Create()
		if err != nil {
			return utils.Error(c, err)
		}
	}

	token := account.GenToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}
