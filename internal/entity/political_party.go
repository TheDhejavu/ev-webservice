package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//PoliticalParty represent the PoliticalParty entity
type PoliticalParty struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Slug      string             `json:"slug" bson:"slug,omitempty"`
	Country   primitive.ObjectID `json:"country" bson:"country,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

type PoliticalPartyRead struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Slug      string             `json:"slug" bson:"slug,omitempty"`
	Country   Country            `json:"country" bson:"country,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

//  PoliticalPartyService represent the PoliticalPartys's usecase
type PoliticalPartyService interface {
	Fetch(ctx context.Context, filter interface{}) (res []PoliticalPartyRead, err error)
	GetByID(ctx context.Context, id string) (PoliticalPartyRead, error)
	GetBySlug(ctx context.Context, slug string) (PoliticalPartyRead, error)
	GetByCountry(ctx context.Context, country string) (PoliticalPartyRead, error)
	Exists(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (bool, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (PoliticalPartyRead, error)
	Store(ctx context.Context, data map[string]interface{}) (PoliticalPartyRead, error)
	Delete(ctx context.Context, id string) error
}

// PoliticalPartyRepository represent the PoliticalPartys's repository contract
type PoliticalPartyRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []PoliticalPartyRead, err error)
	GetByID(ctx context.Context, id string) (PoliticalPartyRead, error)
	GetBySlug(ctx context.Context, slug string) (PoliticalPartyRead, error)
	GetByCountry(ctx context.Context, country string) (PoliticalPartyRead, error)
	GetWithExclude(ctx context.Context, politicalParty map[string]interface{}, exclude map[string]interface{}) (PoliticalParty, error)
	Get(ctx context.Context, filter map[string]interface{}) (PoliticalParty, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (PoliticalPartyRead, error)
	Store(ctx context.Context, politicalParty PoliticalParty) (PoliticalPartyRead, error)
	Delete(ctx context.Context, id string) error
}
