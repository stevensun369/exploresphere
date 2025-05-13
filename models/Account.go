package models

import (
	"backend/db"
	"backend/env"
	"backend/utils"
	"errors"
	"strings"
	"time"

	sj "github.com/brianvoe/sjwt"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type ExchangeSemester struct {
	Year   int    `json:"year" bson:"year"`     // e.g. 2024
	Period string `json:"period" bson:"period"` // e.g. "Winter" or "Summer"
}

func (acc Account) GenToken() (token string) {
	claims, _ := sj.ToClaims(acc)
	claims.SetExpiresAt(time.Now().Add(365 * 24 * time.Hour))

	token = claims.Generate(env.JWT_KEY)
	return
}

func (acc *Account) ParseToken(token string) (err error) {
	verified := sj.Verify(token, env.JWT_KEY)

	if !verified {
		return errors.New("Could not verify token")
	}

	claims, _ := sj.Parse(token)
	err = claims.Validate()
	claims.ToStruct(&acc)

	return
}

func AccountMiddleware(c fiber.Ctx) error {
	var token string

	header := string(c.Get("Authorization"))

	if header != "" &&
		strings.HasPrefix(header, "Bearer") {

		tokens := strings.Fields(header)

		if len(tokens) == 2 {
			token = tokens[1]
		}

		if token == "" {
			return utils.Error(c, errors.New("no token provided"))
		}

		var account Account
		err := account.ParseToken(token)
		if err != nil {
			return utils.Error(c, errors.New("an error has occured"))
		}

		c.Locals("id", account.ID)
		utils.SetLocals(c, "account", account)
	}

	if token == "" {
		return utils.Error(c, errors.New("no token provided"))
	}

	return c.Next()
}

type Account struct {
	ID         string `json:"id" bson:"id"`
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	ProfilePic bool   `json:"profilePic" bson:"profilePic"`

	HomeUniversity        string `json:"homeUniversity" bson:"homeUniversity"`
	DestinationUniversity string `json:"destinationUniversity" bson:"destinationUniversity"`
	Semester              string `json:"semester" bson:"semester"`

	AboutMe string   `json:"aboutMe" bson:"aboutMe"`
	Tags    []string `json:"tags"`

	CreatedAt time.Time `json:"createdAt"`
}

func (acc *Account) Get(id string) error {
	return db.Accounts.FindOne(db.Ctx, bson.M{
		"id": id,
	}).Decode(acc)
}

func (acc *Account) GetByEmail(email string) error {
	return db.Accounts.FindOne(db.Ctx, bson.M{
		"email": email,
	}).Decode(acc)
}

// the heuristic of this function
// is that it already has name, email and hashed password
func (acc *Account) Create() error {

	// generating ID
	acc.ID = utils.GenID(6)
	acc.Tags = []string{}
	acc.CreatedAt = time.Now()

	_, err := db.Accounts.InsertOne(db.Ctx, acc)

	return err
}

func UpdateAccount(id string, updates interface{}) error {
	_, err := db.Accounts.UpdateOne(
		db.Ctx,
		bson.M{"id": id},
		bson.M{
			"$set": updates,
		},
	)

	return err
}

func CheckAccount(email string) (bool, Account) {
	var account Account

	err := db.Accounts.FindOne(
		db.Ctx, bson.M{
			"email": email,
		},
	).Decode(&account)

	if err != nil {
		return false, Account{}
	} else {
		return true, account
	}
}
