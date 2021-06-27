package entity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsensusGroup struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	ServerUrl string             `json:"server_url" bson:"server_url,omitempty"`
	PublicKey string             `json:"public_key" bson:"public_key,omitempty"`
	Country   primitive.ObjectID `json:"country" bson:"country,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type ConsensusGroupRead struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	ServerUrl string             `json:"server_url" bson:"server_url,omitempty"`
	PublicKey string             `json:"public_key" bson:"public_key,omitempty"`
	Country   Country            `json:"country" bson:"country,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

//  ConsensusGroupService represent the ConsensusGroups's usecase
type ConsensusGroupService interface {
	Fetch(ctx context.Context, filter interface{}) (res []ConsensusGroupRead, err error)
	GetByID(ctx context.Context, id string) (ConsensusGroupRead, error)
	GetByPubKey(ctx context.Context, publicKey string) (ConsensusGroupRead, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (ConsensusGroupRead, error)
	Create(ctx context.Context, data map[string]interface{}) (ConsensusGroupRead, error)
	Exists(ctx context.Context, group map[string]interface{}, exclude map[string]interface{}) (bool, error)
	Delete(ctx context.Context, id string) error
}

// ConsensusGroupRepository represent the ConsensusGroups's repository contract
type ConsensusGroupRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []ConsensusGroupRead, err error)
	GetByID(ctx context.Context, id string) (ConsensusGroupRead, error)
	GetByPubKey(ctx context.Context, publicKey string) (ConsensusGroupRead, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (ConsensusGroupRead, error)
	Create(ctx context.Context, group ConsensusGroup) (ConsensusGroupRead, error)
	Get(ctx context.Context, filter map[string]interface{}) (ConsensusGroupRead, error)
	GetWithExclude(ctx context.Context, group map[string]interface{}, exclude map[string]interface{}) (ConsensusGroup, error)
	Delete(ctx context.Context, id string) error
}
