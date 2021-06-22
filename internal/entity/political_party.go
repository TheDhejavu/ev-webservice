package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PoliticalParty struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Slug      string             `json:"slug" bson:"slug,omitempty"`
	CountryId string             `json:"country_id" bson:"country_id,omitempty"`
	Country   primitive.ObjectID `json:"country" bson:"country,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

//  PoliticalPartyService represent the PoliticalPartys's usecase
type PoliticalPartyService interface {
	Fetch(ctx context.Context, filter interface{}) (res []PoliticalParty, err error)
	GetByID(ctx context.Context, id string) (PoliticalParty, error)
	GetBySlug(ctx context.Context, slug string) (PoliticalParty, error)
	GetByCountry(ctx context.Context, country string) (PoliticalParty, error)
	Update(ctx context.Context, id string, data interface{}) (PoliticalParty, error)
	Store(ctx context.Context, PoliticalParty PoliticalParty) (PoliticalParty, error)
	Delete(ctx context.Context, id string) error
}

// PoliticalPartyRepository represent the PoliticalPartys's repository contract
type PoliticalPartyRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []PoliticalParty, err error)
	GetByID(ctx context.Context, id string) (PoliticalParty, error)
	GetBySlug(ctx context.Context, slug string) (PoliticalParty, error)
	GetByCountry(ctx context.Context, country string) (PoliticalParty, error)
	Update(ctx context.Context, id string, data interface{}) (PoliticalParty, error)
	Store(ctx context.Context, PoliticalParty PoliticalParty) (PoliticalParty, error)
	Delete(ctx context.Context, id string) error
}
