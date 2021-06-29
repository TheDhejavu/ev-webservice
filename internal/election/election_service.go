package election

import (
	"context"
	"encoding/json"

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

func (service *electionService) Fetch(ctx context.Context, filter interface{}) (res []entity.ElectionRead, err error) {
	res, err = service.electionRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}
	return
}
func (service *electionService) GetByID(ctx context.Context, id string) (res entity.ElectionRead, err error) {
	res, err = service.electionRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}
func (service *electionService) Update(ctx context.Context, id string, data map[string]interface{}) (res entity.ElectionRead, err error) {
	res, err = service.electionRepo.Update(ctx, id, data)
	if err != nil {
		return
	}
	return
}
func (service *electionService) Create(ctx context.Context, data map[string]interface{}) (res entity.ElectionRead, err error) {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		return
	}

	election := &entity.Election{}
	if err = json.Unmarshal(jsonbody, &election); err != nil {
		return
	}

	res, err = service.electionRepo.Create(ctx, *election)
	if err != nil {
		return
	}
	return
}
func (service *electionService) Delete(ctx context.Context, id string) (err error) {
	value, _ := service.Exists(ctx, map[string]interface{}{"_id": id}, nil)
	if value == false {
		return entity.ErrNotFound
	}

	err = service.electionRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}

func (service *electionService) Exists(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (res bool, err error) {
	if exclude == nil {
		_, err = service.electionRepo.Get(ctx, filter)
	} else {
		_, err = service.electionRepo.GetWithExclude(ctx, filter, exclude)
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

func (services *electionService) GetResult(ctx context.Context, filter map[string]interface{}) (res []entity.ElectionRead, err error) {
	return
}
