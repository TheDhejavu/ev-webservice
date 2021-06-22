package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Identity represent the Identity entity
type Identity struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Digit            int64              `json:"digit" bson:"digit,omitempty"`
	FirstName        string             `json:"first_name" bson:"first_name,omitempty"`
	LastName         string             `json:"last_name" bson:"last_name,omitempty"`
	Country          primitive.ObjectID `json:"country" bson:"country,omitempty"`
	State            State              `json:"state" bson:"state,omitempty"`
	Residence        Residence          `json:"residence" bson:"residence,omitempty"`
	Email            string             `json:"email" bson:"email"`
	Password         string             `json:"password" bson:"password"`
	BirthCertificate string             `json:"birth_certificate" bson:"birth_certificate"`
	NINCard          string             `json:"nin_card" bson:"nin_card"`
	VoterCard        string             `json:"voter_card" bson:"voter_card"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type State struct {
	City    string `json:"city" bson:"city,omitempty"`
	Address string `json:"address" bson:"address,omitempty"`
}

type Residence struct {
	Country string `json:"country" bson:"city,omitempty"`
	City    string `json:"city" bson:"city,omitempty"`
	Address string `json:"address" bson:"address,omitempty"`
}

//  IdentityService represent the Identitys's usecase
type IdentityService interface {
	Fetch(ctx context.Context, filter interface{}) (res []Identity, err error)
	GetByID(ctx context.Context, id string) (Identity, error)
	GetByDigit(ctx context.Context, digit int64) (Identity, error)
	CheckEmailIsTaken(ctx context.Context, email string) (Identity, error)
	GetByEmail(ctx context.Context, email string) (Identity, error)
	Update(ctx context.Context, id string, data interface{}) (Identity, error)
	Create(ctx context.Context, Identity Identity) (Identity, error)
	Delete(ctx context.Context, id string) error
}

// IdentityRepository represent the Identitys's repository contract
type IdentityRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []Identity, err error)
	GetByID(ctx context.Context, id string) (Identity, error)
	GetByDigit(ctx context.Context, digit int64) (Identity, error)
	GetByEmail(ctx context.Context, email string) (Identity, error)
	Update(ctx context.Context, id string, data interface{}) (Identity, error)
	Create(ctx context.Context, Identity Identity) (Identity, error)
	Delete(ctx context.Context, id string) error
}
