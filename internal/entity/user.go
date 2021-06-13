package entity

import (
	"time"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

//User represent the user entity
type User struct {
	// ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string    `json:"username" bson:"username,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Status    int       `json:"status" bson:"status"`
	Roles     string    `json:"roles" bson:"roles"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty""`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty""`
}
