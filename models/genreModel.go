package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Genre struct {
	Genre_ID   primitive.ObjectID `bson:"_id"`
	Name       *string            `json:"name" validate:"required,min=4,max=100"`
	Creator_id *string            `json:"creator_id" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Genre_id   string             `json:"genre_id"`
}
