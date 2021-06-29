package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Country struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Flag      string             `json:"flag,omitempty" bson:"flag,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug      string             `json:"slug,omitempty" bson:"slug,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//  CountryService represent the Countrys's usecase
type CountryService interface {
	Fetch(ctx context.Context, filter interface{}) (res []Country, err error)
	GetByID(ctx context.Context, id string) (Country, error)
	IdExists(ctx context.Context, id string) (bool, error)
	GetBySlug(ctx context.Context, slug string) (Country, error)
	SlugExists(ctx context.Context, slug string, id interface{}) (bool, error)
	NameExists(ctx context.Context, name string, id interface{}) (bool, error)
	GetByName(ctx context.Context, name string) (Country, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (Country, error)
	Store(ctx context.Context, country map[string]interface{}) (Country, error)
	Delete(ctx context.Context, id string) error
}

// CountryRepository represent the Countrys's repository contract
type CountryRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []Country, err error)
	GetByID(ctx context.Context, id string) (Country, error)
	GetBySlug(ctx context.Context, slug string) (Country, error)
	GetByName(ctx context.Context, name string) (Country, error)
	GetWithExclude(ctx context.Context, country map[string]interface{}, exclude map[string]interface{}) (Country, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (Country, error)
	Store(ctx context.Context, Country Country) (Country, error)
	Delete(ctx context.Context, id string) error
}
