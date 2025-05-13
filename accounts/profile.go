package accounts

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func postProfile(c fiber.Ctx) error {
	var body struct {
		Name                  string   `json:"name"`
		HomeUniversity        string   `json:"homeUniversity"`
		DestinationUniversity string   `json:"destinationUniversity"`
		Semester              string   `json:"semester"`
		AboutMe               string   `json:"aboutMe"`
		Tags                  []string `json:"tags"`
	}
	json.Unmarshal(c.Body(), &body)

	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	fmt.Println(body.Tags)

	err := models.UpdateAccount(
		account.ID,
		bson.M{
			"name":                  body.Name,
			"homeUniversity":        body.HomeUniversity,
			"destinationUniversity": body.DestinationUniversity,
			"semester":              body.Semester,
			"aboutMe":               body.AboutMe,
			"tags":                  body.Tags,
		},
	)
	if err != nil {
		return utils.Error(c, err)
	}

	account.Name = body.Name
	account.HomeUniversity = body.HomeUniversity
	account.DestinationUniversity = body.DestinationUniversity
	account.Semester = body.Semester
	account.AboutMe = body.AboutMe
	account.Tags = body.Tags

	token := account.GenToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}

func postProfilePic(c fiber.Ctx) error {
	account := models.Account{}
	utils.GetLocals(c, "account", &account)

	file, err := c.FormFile("image")
	if err != nil {
		return utils.Error(c, err)
	}

	savePath := fmt.Sprintf("./profilepics/%s", account.ID)
	if err := c.SaveFile(file, savePath); err != nil {
		return utils.Error(c, err)
	}

	account.ProfilePic = true

	err = models.UpdateAccount(account.ID, bson.M{
		"profilePic": true,
	})
	if err != nil {
		return utils.Error(c, err)
	}

	token := account.GenToken()

	return c.JSON(bson.M{
		"account": account,
		"token":   token,
	})
}
