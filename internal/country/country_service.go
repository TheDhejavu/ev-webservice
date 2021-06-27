package country

import (
	"context"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type countryService struct {
	countryRepo entity.CountryRepository
	logger      log.Logger
}

// NewCountryService creates a new country service.
func NewCountryService(countryRepo entity.CountryRepository, logger log.Logger) entity.CountryService {
	return &countryService{
		countryRepo: countryRepo,
		logger:      logger,
	}
}

// Fetch returns the countries with the specified filter.
func (service *countryService) Fetch(ctx context.Context, filter interface{}) (res []entity.Country, err error) {
	res, err = service.countryRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}
	return
}

// GetById returns the country with the specified country ID.
func (service *countryService) GetByID(ctx context.Context, id string) (res entity.Country, err error) {
	res, err = service.countryRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

// GetByName returns the country with the specified country ID.
func (service *countryService) GetByName(ctx context.Context, name string) (res entity.Country, err error) {
	res, err = service.countryRepo.GetByName(ctx, name)
	if err != nil {
		return
	}
	return
}

// GetBySlug returns the country with the specified country slug.
func (service *countryService) GetBySlug(ctx context.Context, slug string) (res entity.Country, err error) {
	res, err = service.countryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return
	}
	return
}

// Update updates the country with the specified ID.
func (service *countryService) Update(ctx context.Context, id string, data map[string]interface{}) (res entity.Country, err error) {
	res, err = service.countryRepo.Update(ctx, id, data)
	if err != nil {
		return
	}
	return
}

// Store stores new country
func (service *countryService) Store(ctx context.Context, country map[string]interface{}) (res entity.Country, err error) {
	new_country := entity.Country{
		Flag: fmt.Sprintf("%v", country["flag"]),
		Slug: fmt.Sprintf("%v", country["slug"]),
		Name: fmt.Sprintf("%v", country["name"]),
	}
	res, err = service.countryRepo.Store(ctx, new_country)
	if err != nil {
		return
	}
	return
}

// IdExsits checks if a country exists with the specified ID.
func (service *countryService) IdExists(ctx context.Context, id string) (bool, error) {
	_, err := service.countryRepo.GetByID(ctx, id)

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, err
}

// Delete deletes the user with the specified ID.
func (service *countryService) Delete(ctx context.Context, id string) (err error) {
	value, _ := service.IdExists(ctx, id)
	if value == false {
		return entity.ErrNotFound
	}

	err = service.countryRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}

// NameExists checks if specified name exist already.
func (service *countryService) NameExists(ctx context.Context, name string, id interface{}) (bool, error) {
	var err error
	if id != nil {
		country := map[string]interface{}{
			"name": name,
		}
		exlude := map[string]interface{}{
			"id": id,
		}

		_, err = service.countryRepo.GetWithExclude(ctx, country, exlude)
	} else {
		_, err = service.countryRepo.GetByName(ctx, name)
	}
	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, err
}

// SlugExists checks if specified slug exist already.
func (service *countryService) SlugExists(ctx context.Context, slug string, id interface{}) (bool, error) {
	var err error
	if id != nil {
		country := map[string]interface{}{
			"slug": slug,
		}
		exlude := map[string]interface{}{
			"id": id,
		}
		_, err = service.countryRepo.GetWithExclude(ctx, country, exlude)
	} else {
		_, err = service.countryRepo.GetBySlug(ctx, slug)
	}

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, err
}
