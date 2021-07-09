package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Identity represent the Identity entity
type Identity struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Digits           uint64             `json:"digits" bson:"digits,omitempty"`
	FirstName        string             `json:"first_name" bson:"first_name,omitempty"`
	LastName         string             `json:"last_name" bson:"last_name,omitempty"`
	Origin           Origin             `json:"origin" bson:"origin,omitempty"`
	Residence        Residence          `json:"residence" bson:"residence,omitempty"`
	Email            string             `json:"email" bson:"email"`
	Password         string             `json:"password" bson:"password"`
	BirthCertificate string             `json:"birth_certificate" bson:"birth_certificate"`
	NationalIdCard   string             `json:"national_id_card" bson:"national_id_card"`
	VoterCard        string             `json:"voter_card" bson:"voter_card"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at,omitempty"`
}
type Origin struct {
	Country primitive.ObjectID `json:"country" bson:"country,omitempty"`
	State   string             `json:"state" bson:"state,omitempty"`
	City    string             `json:"city" bson:"city,omitempty"`
	Address string             `json:"address" bson:"address,omitempty"`
}

type Residence struct {
	Country primitive.ObjectID `json:"country" bson:"country,omitempty"`
	City    string             `json:"city" bson:"city,omitempty"`
	State   string             `json:"state" bson:"state,omitempty"`
	Address string             `json:"address" bson:"address,omitempty"`
}

// Identity represent the Identity entity
type IdentityRead struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Digits           uint64             `json:"digits" bson:"digits,omitempty"`
	FirstName        string             `json:"first_name" bson:"first_name,omitempty"`
	LastName         string             `json:"last_name" bson:"last_name,omitempty"`
	Origin           OriginRead         `json:"origin" bson:"origin,omitempty"`
	Residence        ResidenceRead      `json:"residence" bson:"residence,omitempty"`
	Email            string             `json:"email" bson:"email,omitempty"`
	Password         string             `json:"password,omitempty" bson:"password"`
	BirthCertificate string             `json:"birth_certificate" bson:"birth_certificate"`
	NationalIdCard   string             `json:"national_id_card" bson:"national_id_card"`
	VoterCard        string             `json:"voter_card" bson:"voter_card"`
	Wallet           struct {
		PublicMainKey string `json:"public_main_key"`
		PublicViewKey string `json:"public_view_key"`
		Certificate   string `json:"certificate"`
	} `json:"wallet"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
}

type OriginRead struct {
	Country Country `json:"country" bson:"country,omitempty"`
	State   string  `json:"state" bson:"state,omitempty"`
	City    string  `json:"city" bson:"city,omitempty"`
	Address string  `json:"address" bson:"address,omitempty"`
}

type ResidenceRead struct {
	Country Country `json:"country" bson:"country,omitempty"`
	City    string  `json:"city" bson:"city,omitempty"`
	State   string  `json:"state" bson:"state,omitempty"`
	Address string  `json:"address" bson:"address,omitempty"`
}

//  IdentityService represent the Identitys's usecase
type IdentityService interface {
	Fetch(ctx context.Context, filter interface{}) (res []IdentityRead, err error)
	GetByID(ctx context.Context, id string) (IdentityRead, error)
	GetByDigits(ctx context.Context, digits uint64) (IdentityRead, error)
	GetByEmail(ctx context.Context, email string) (IdentityRead, error)
	Exists(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (bool, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (IdentityRead, error)
	Create(ctx context.Context, identity map[string]interface{}, imageTmpPaths []string) (IdentityRead, error)
	Delete(ctx context.Context, id string) error
}

// IdentityRepository represent the Identitys's repository contract
type IdentityRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []IdentityRead, err error)
	GetByID(ctx context.Context, id string) (IdentityRead, error)
	GetByEmail(ctx context.Context, slug string) (IdentityRead, error)
	GetByDigits(ctx context.Context, digits uint64) (IdentityRead, error)
	GetByCountry(ctx context.Context, country string) (IdentityRead, error)
	GetWithExclude(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (Identity, error)
	Get(ctx context.Context, filter map[string]interface{}) (Identity, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (IdentityRead, error)
	Create(ctx context.Context, identity Identity) (IdentityRead, error)
	Delete(ctx context.Context, id string) error
}
