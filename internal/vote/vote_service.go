package vote

import (
	"context"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type voteService struct {
	logger log.Logger
}

func NewvoteService(logger log.Logger) entity.VoteService {
	return &voteService{
		logger: logger,
	}
}

func (vs *voteService) Fetch(ctx context.Context, filter interface{}) (res []entity.Vote, err error) {
	return
}
func (vs *voteService) GetByID(ctx context.Context, id string) (res entity.Vote, err error) {
	return
}
func (vs *voteService) Cast(ctx context.Context, vote entity.Vote) (res entity.Vote, err error) {
	return
}
func (vs *voteService) Start(ctx context.Context, id string) (res entity.Vote, err error) {
	return
}
func (vs *voteService) Stop(ctx context.Context, id string) (res entity.Vote, err error) {
	return
}
