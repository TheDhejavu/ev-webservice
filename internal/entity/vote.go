package entity

import (
	"context"
	"time"
)

//Vote represent the Vote entity
type Vote struct {
	ID             []byte    `json:"id" bson:"_id,omitempty"`
	PoliticalParty string    `json:"political_party" bson:"political_party,omitempty"`
	Election       Election  `json:"election" bson:"election,omitempty"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at,omitempty"`
}

//  VoteService represent the Votes's usecase
type VoteService interface {
	Fetch(ctx context.Context, filter interface{}) (res []Vote, err error)
	GetByID(ctx context.Context, id string) (Vote, error)
	Cast(ctx context.Context, Vote Vote) (Vote, error)
	Start(ctx context.Context, id string) (Vote, error)
	Stop(ctx context.Context, id string) (Vote, error)
}
