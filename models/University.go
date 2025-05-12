package models

type University struct {
	ID      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	City    string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
}
