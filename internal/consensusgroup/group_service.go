package consensusgroup

import (
	"context"
	"encoding/json"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type groupService struct {
	groupRepo entity.ConsensusGroupRepository
	logger    log.Logger
}

func NewGroupService(ConsensusGroupRepo entity.ConsensusGroupRepository, logger log.Logger) entity.ConsensusGroupService {
	return &groupService{
		groupRepo: ConsensusGroupRepo,
		logger:    logger,
	}
}

func (service *groupService) Fetch(ctx context.Context, filter interface{}) (res []entity.ConsensusGroupRead, err error) {
	res, err = service.groupRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}
	return
}
func (service *groupService) GetByID(ctx context.Context, id string) (res entity.ConsensusGroupRead, err error) {
	res, err = service.groupRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}
func (service *groupService) Update(ctx context.Context, id string, data map[string]interface{}) (res entity.ConsensusGroupRead, err error) {
	res, err = service.groupRepo.Update(ctx, id, data)
	if err != nil {
		return
	}
	return
}
func (service *groupService) Create(ctx context.Context, group map[string]interface{}) (res entity.ConsensusGroupRead, err error) {
	jsonbody, err := json.Marshal(group)
	if err != nil {
		return
	}

	newGroup := &entity.ConsensusGroup{}
	if err = json.Unmarshal(jsonbody, &newGroup); err != nil {
		return
	}

	res, err = service.groupRepo.Create(ctx, *newGroup)
	if err != nil {
		return
	}
	return
}
func (service *groupService) Delete(ctx context.Context, id string) (err error) {
	value, _ := service.Exists(ctx, map[string]interface{}{"_id": id}, nil)

	if value == false {
		return entity.ErrNotFound
	}

	err = service.groupRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}

func (service *groupService) GetByPubKey(ctx context.Context, public_key string) (res entity.ConsensusGroupRead, err error) {
	res, err = service.groupRepo.GetByPubKey(ctx, public_key)
	if err != nil {
		return
	}
	return
}

func (service *groupService) Exists(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (res bool, err error) {
	if exclude == nil {
		_, err = service.groupRepo.Get(ctx, filter)
	} else {
		_, err = service.groupRepo.GetWithExclude(ctx, filter, exclude)
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
