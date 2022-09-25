package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       *string            `json:"name" validate:"required"`
	Topic      *string            `json:"topic" validate:"required"`
	Genre_id   *string            `json:"genre_id" validate:"required"`
	Movie_URL  *string            `json:"movie_url" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
