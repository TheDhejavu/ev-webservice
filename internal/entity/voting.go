package entity

import (
	"context"
)

//Vote represent the Vote entity
type Vote struct {
	ID             []byte   `json:"id" bson:"_id,omitempty"`
	PoliticalParty string   `json:"political_party" bson:"political_party,omitempty"`
	Election       Election `json:"election" bson:"election,omitempty"`
}

//  VoteService represent the Votes's usecase
type VotingService interface {
	GetResults(ctx context.Context, id string) (res []*CandidateRead, err error)
	CastVote(ctx context.Context, userId, electionId, candidate string) (res *CandidateRead, err error)
	Start(ctx context.Context, id string) (Vote, error)
	Stop(ctx context.Context, id string) (Vote, error)
}
