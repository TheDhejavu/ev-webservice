package political_party

import (
	"context"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type politicalPartyService struct {
	PoliticalPartyRepo entity.PoliticalPartyRepository
	logger             log.Logger
}

func NewpoliticalPartyService(PoliticalPartyRepo entity.PoliticalPartyRepository, logger log.Logger) entity.PoliticalPartyService {
	return &politicalPartyService{
		PoliticalPartyRepo: PoliticalPartyRepo,
		logger:             logger,
	}
}

func (p *politicalPartyService) Fetch(ctx context.Context, filter interface{}) (res []entity.PoliticalParty, err error) {
	return
}
func (p *politicalPartyService) GetByID(ctx context.Context, id string) (res entity.PoliticalParty, err error) {
	return
}
func (p *politicalPartyService) Update(ctx context.Context, id string, data interface{}) (res entity.PoliticalParty, err error) {
	return
}
func (p *politicalPartyService) Store(ctx context.Context, PoliticalParty entity.PoliticalParty) (res entity.PoliticalParty, err error) {
	return
}
func (p *politicalPartyService) Delete(ctx context.Context, id string) error {
	return nil
}

func (p *politicalPartyService) GetBySlug(ctx context.Context, slug string) (res entity.PoliticalParty, err error) {
	return
}
func (p *politicalPartyService) GetByCountry(ctx context.Context, country string) (res entity.PoliticalParty, err error) {
	return
}
