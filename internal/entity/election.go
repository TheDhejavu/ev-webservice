package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Election represent the Election entity
type Election struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title           string             `json:"title" bson:"title,omitempty"`
	Description     string             `json:"description" bson:"description,omitempty"`
	Country         primitive.ObjectID `json:"country" bson:"country,omitempty"`
	Candidates      []Candidate        `json:"candidates" bson:"candidates,omitempty"`
	BlockTxRef      string             `json:"block_tx_reference" bson:"tx_refrence,omitempty"`
	AccreditationAt AccreditationAt    `json:"accrediation_at" bson:"accrediation_at,omitempty"`
	VoteAt          VoteAt             `json:"vote_at" bson:"vote_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type Candidate struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FullName       string             `json:"title" bson:"title,omitempty"`
	Position       string             `json:"position" bson:"position,omitempty"`
	PoliticalParty primitive.ObjectID `json:"political_party" bson:"political_party,omitempty"`
}
type AccreditationAt struct {
	Start time.Time `json:"start" bson:"start,omitempty"`
	End   time.Time `json:"end" bson:"end,omitempty"`
}

type VoteAt struct {
	Start time.Time `json:"start" bson:"start,omitempty"`
	End   time.Time `json:"end" bson:"end,omitempty"`
}

// Election Read
type ElectionRead struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title            string             `json:"title" bson:"title,omitempty"`
	Description      string             `json:"description" bson:"description,omitempty"`
	AccreditationAt  AccreditationAt    `json:"accrediation_at" bson:"accrediation_at,omitempty"`
	VoteAt           VoteAt             `json:"vote_at" bson:"vote_at,omitempty"`
	Country          Country            `json:"country" bson:"country,omitempty"`
	Candidates       []CandidateRead    `json:"candidates" bson:"candidates,omitempty"`
	BlockTxReference []byte             `json:"block_tx_reference" bson:"tx_refrence,omitempty"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type CandidateRead struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FullName       string             `json:"title" bson:"title,omitempty"`
	Position       string             `json:"position" bson:"position,omitempty"`
	PoliticalParty PoliticalParty     `json:"political_party" bson:"political_party,omitempty"`
}

//  ElectionService represent the Elections's usecase
type ElectionService interface {
	Fetch(ctx context.Context, filter interface{}) (res []Election, err error)
	GetByID(ctx context.Context, id string) (Election, error)
	Update(ctx context.Context, id string, data interface{}) (Election, error)
	Create(ctx context.Context, election Election) (Election, error)
	Delete(ctx context.Context, id string) error
	GetResult(ctx context.Context, filter interface{}) (res []Election, err error)
}

// ElectionRepository represent the Elections's repository contract
type ElectionRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []Election, err error)
	GetByID(ctx context.Context, id string) (Election, error)
	Update(ctx context.Context, id string, data interface{}) (Election, error)
	Create(ctx context.Context, Election Election) (Election, error)
	Delete(ctx context.Context, id string) error
}
