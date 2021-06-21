package entity

import "context"

type Country struct {
	Flag string `json:"flag" bson:"flag,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Slug string `json:"slug" bson:"slug,omitempty"`
}

//  CountryService represent the Countrys's usecase
type CountryService interface {
	Fetch(ctx context.Context, filter interface{}) (res []Country, err error)
	GetByID(ctx context.Context, id string) (Country, error)
	Update(ctx context.Context, id string, data interface{}) (Country, error)
	Store(ctx context.Context, Country Country) (Country, error)
	Delete(ctx context.Context, id string) error
}

// CountryRepository represent the Countrys's repository contract
type CountryRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []Country, err error)
	GetByID(ctx context.Context, id string) (Country, error)
	Update(ctx context.Context, id string, data interface{}) (Country, error)
	Store(ctx context.Context, Country Country) (Country, error)
	Delete(ctx context.Context, id string) error
}