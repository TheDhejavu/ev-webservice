package consensus_group

import (
	"context"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type ConsensusGroupService struct {
	ConsensusGroupRepo entity.ConsensusGroupRepository
	logger             log.Logger
}

func NewConsensusGroupService(ConsensusGroupRepo entity.ConsensusGroupRepository, logger log.Logger) entity.ConsensusGroupService {
	return &ConsensusGroupService{
		ConsensusGroupRepo: ConsensusGroupRepo,
		logger:             logger,
	}
}

func (group *ConsensusGroupService) Fetch(ctx context.Context, filter interface{}) (res []entity.ConsensusGroup, err error) {
	return
}
func (group *ConsensusGroupService) GetByID(ctx context.Context, id string) (res entity.ConsensusGroup, err error) {
	return
}
func (group *ConsensusGroupService) Update(ctx context.Context, id string, data interface{}) (res entity.ConsensusGroup, err error) {
	return
}
func (group *ConsensusGroupService) Create(ctx context.Context, ConsensusGroup entity.ConsensusGroup) (res entity.ConsensusGroup, err error) {
	return
}
func (group *ConsensusGroupService) Delete(ctx context.Context, id string) error {
	return nil
}

func (group *ConsensusGroupService) GetByPubKey(ctx context.Context, public_key []byte) (res entity.ConsensusGroup, err error) {
	return
}
