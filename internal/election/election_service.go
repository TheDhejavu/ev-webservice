package election

import (
	"context"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type electionService struct {
	electionRepo entity.ElectionRepository
	logger       log.Logger
}

func NewElectionService(electionRepo entity.ElectionRepository, logger log.Logger) entity.ElectionService {
	return &electionService{
		electionRepo: electionRepo,
		logger:       logger,
	}
}

func (es *electionService) Fetch(ctx context.Context, filter interface{}) (res []entity.Election, err error) {
	return
}
func (es *electionService) GetByID(ctx context.Context, id string) (res entity.Election, err error) {
	return
}
func (es *electionService) Update(ctx context.Context, id string, data interface{}) (res entity.Election, err error) {
	return
}
func (es *electionService) Create(ctx context.Context, Election entity.Election) (res entity.Election, err error) {
	return
}
func (es *electionService) Delete(ctx context.Context, id string) error {
	return nil
}
func (es *electionService) GetResult(ctx context.Context, filter interface{}) (res []entity.Election, err error) {
	return
}
