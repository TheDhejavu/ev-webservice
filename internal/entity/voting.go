package entity

import (
	"context"
)

//Vote represent the Vote entity
type Vote struct {
	Election ElectionRead `json:"election" bson:"election,omitempty"`
}

//  VoteService represent the Votes's usecase
type VotingService interface {
	CastVote(ctx context.Context, userId, electionId, candidate string) (res *CandidateRead, err error)
	Start(ctx context.Context, id string) (Vote, error)
	Stop(ctx context.Context, id string) (Vote, error)
}
