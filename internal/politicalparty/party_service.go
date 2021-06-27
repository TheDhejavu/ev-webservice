package politicalparty

import (
	"context"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type politicalPartyService struct {
	politicalPartyRepo entity.PoliticalPartyRepository
	logger             log.Logger
}

func NewPoliticalPartyService(PoliticalPartyRepo entity.PoliticalPartyRepository, logger log.Logger) entity.PoliticalPartyService {
	return &politicalPartyService{
		politicalPartyRepo: PoliticalPartyRepo,
		logger:             logger,
	}
}

func (service *politicalPartyService) Fetch(ctx context.Context, filter interface{}) (res []entity.PoliticalPartyRead, err error) {
	res, err = service.politicalPartyRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}
	return
}
func (service *politicalPartyService) GetByID(ctx context.Context, id string) (res entity.PoliticalPartyRead, err error) {
	res, err = service.politicalPartyRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}
func (service *politicalPartyService) Update(ctx context.Context, id string, data map[string]interface{}) (res entity.PoliticalPartyRead, err error) {
	fmt.Println(id)
	res, err = service.politicalPartyRepo.Update(ctx, id, data)
	if err != nil {
		return
	}
	return
}
func (service *politicalPartyService) Store(ctx context.Context, politicalParty map[string]interface{}) (res entity.PoliticalPartyRead, err error) {
	country := fmt.Sprintf("%v", politicalParty["country"])
	_id, err := primitive.ObjectIDFromHex(country)
	new_party := entity.PoliticalParty{
		Name:    fmt.Sprintf("%v", politicalParty["name"]),
		Slug:    fmt.Sprintf("%v", politicalParty["slug"]),
		Country: _id,
	}
	res, err = service.politicalPartyRepo.Store(ctx, new_party)
	if err != nil {
		return
	}
	return
}

// IdExsits checks if a country exists with the specified ID.
func (service *politicalPartyService) IdExists(ctx context.Context, id string) (bool, error) {
	_, err := service.politicalPartyRepo.GetByID(ctx, id)

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

func (service *politicalPartyService) Delete(ctx context.Context, id string) (err error) {
	value, _ := service.IdExists(ctx, id)
	if value == false {
		return entity.ErrNotFound
	}

	err = service.politicalPartyRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}

func (service *politicalPartyService) GetBySlug(ctx context.Context, slug string) (res entity.PoliticalPartyRead, err error) {
	res, err = service.politicalPartyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return
	}
	return
}
func (service *politicalPartyService) GetByCountry(ctx context.Context, country string) (res entity.PoliticalPartyRead, err error) {
	res, err = service.politicalPartyRepo.GetByCountry(ctx, country)
	if err != nil {
		return
	}
	return
}

// Exists checks if specified party exist already.
func (service *politicalPartyService) Exists(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (bool, error) {
	var err error
	if exclude == nil {
		_, err = service.politicalPartyRepo.Get(ctx, filter)
	} else {
		_, err = service.politicalPartyRepo.GetWithExclude(ctx, filter, exclude)
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
