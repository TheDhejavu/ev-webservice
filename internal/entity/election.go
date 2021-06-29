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
	Phase           string             `json:"phase" bson:"phase,omitempty"`
	Country         primitive.ObjectID `json:"country" bson:"country,omitempty"`
	Candidates      []Candidate        `json:"candidates" bson:"candidates,omitempty"`
	AccreditationAt ElectionAt         `json:"accrediation_at,omitempty" bson:"accrediation_at,omitempty"`
	VoteAt          ElectionAt         `json:"vote_at,omitempty" bson:"vote_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// Election Read
type ElectionRead struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title           string             `json:"title,omitempty" bson:"title,omitempty"`
	Description     string             `json:"description,omitempty" bson:"description,omitempty"`
	Phase           string             `json:"phase" bson:"phase,omitempty"`
	AccreditationAt ElectionAt         `json:"accrediation_at,omitempty" bson:"accrediation_at,omitempty"`
	VoteAt          ElectionAt         `json:"vote_at,omitempty" bson:"vote_at,omitempty"`
	Country         Country            `json:"country,omitempty" bson:"country,omitempty"`
	Candidates      []CandidateRead    `json:"candidates,omitempty" bson:"candidates,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//Candidate
type Candidate struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FullName       string             `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Position       string             `json:"position,omitempty" bson:"position,omitempty"`
	PoliticalParty primitive.ObjectID `json:"political_party,omitempty" bson:"political_party,omitempty"`
}

type CandidateRead struct {
	ID             string         `json:"id,omitempty" bson:"_id,omitempty"`
	FullName       string         `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Position       string         `json:"position,omitempty" bson:"position,omitempty"`
	PoliticalParty PoliticalParty `json:"political_party,omitempty" bson:"political_party,omitempty"`
}
type ElectionAt struct {
	TxStartRef string    `json:"tx_start_ref" bson:"tx_start_ref,omitempty"`
	TxEndRef   string    `json:"tx_end_ref" bson:"tx_send_ref,omitempty"`
	Start      time.Time `json:"start,omitempty" bson:"start,omitempty"`
	End        time.Time `json:"end,omitempty" bson:"end,omitempty"`
}

//  ElectionService represent the Elections's usecase
type ElectionService interface {
	Fetch(ctx context.Context, filter interface{}) (res []ElectionRead, err error)
	GetByID(ctx context.Context, id string) (ElectionRead, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (ElectionRead, error)
	Create(ctx context.Context, election map[string]interface{}) (ElectionRead, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, election map[string]interface{}, exclude map[string]interface{}) (bool, error)
	GetResult(ctx context.Context, filter map[string]interface{}) (res []ElectionRead, err error)
}

// ElectionRepository represent the Elections's repository contract
type ElectionRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []ElectionRead, err error)
	GetByID(ctx context.Context, id string) (ElectionRead, error)
	GetBySlug(ctx context.Context, slug string) (ElectionRead, error)
	GetByCountry(ctx context.Context, country string) (ElectionRead, error)
	GetWithExclude(ctx context.Context, Election map[string]interface{}, exclude map[string]interface{}) (Election, error)
	Get(ctx context.Context, filter map[string]interface{}) (Election, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (ElectionRead, error)
	Create(ctx context.Context, election Election) (ElectionRead, error)
	Delete(ctx context.Context, id string) error
}
