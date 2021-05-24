package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct to store user data
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	DeletedAt time.Time          `json:"-" bson:"deleted_at,omitempty"`
}
