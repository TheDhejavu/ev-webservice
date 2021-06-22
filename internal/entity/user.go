package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User represent the user entity
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Role      string             `json:"roles" bson:"roles"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

//  UserService represent the Users's usecase
type UserService interface {
	Fetch(ctx context.Context, filter interface{}) (res []User, err error)
	GetByID(ctx context.Context, id string) (User, error)
	CheckEmailIsTaken(ctx context.Context, email string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, id string, data interface{}) (User, error)
	Store(ctx context.Context, User User) (User, error)
	Delete(ctx context.Context, id string) error
}

// UserRepository represent the Users's repository contract
type UserRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []User, err error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, id string, data interface{}) (User, error)
	Store(ctx context.Context, User User) (User, error)
	Delete(ctx context.Context, id string) error
}
