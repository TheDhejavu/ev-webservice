package country

import (
	"context"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type countryService struct {
	countryRepo entity.CountryRepository
	logger      log.Logger
}

func NewCountryService(countryRepo entity.CountryRepository, logger log.Logger) entity.CountryService {
	return &countryService{
		countryRepo: countryRepo,
		logger:      logger,
	}
}

func (es *countryService) Fetch(ctx context.Context, filter interface{}) (res []entity.Country, err error) {
	return
}
func (es *countryService) GetByID(ctx context.Context, id string) (res entity.Country, err error) {
	return
}
func (es *countryService) Update(ctx context.Context, id string, data interface{}) (res entity.Country, err error) {
	return
}
func (es *countryService) Store(ctx context.Context, country entity.Country) (res entity.Country, err error) {
	return
}
func (es *countryService) Delete(ctx context.Context, id string) error {
	return nil
}
