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
	PublicKey []byte             `json:"public_key" bson:"public_key,omitempty"`
	Country   primitive.ObjectID `json:"country" bson:"country,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

//  ConsensusGroupService represent the ConsensusGroups's usecase
type ConsensusGroupService interface {
	Fetch(ctx context.Context, filter interface{}) (res []ConsensusGroup, err error)
	GetByID(ctx context.Context, id string) (ConsensusGroup, error)
	GetByPubKey(ctx context.Context, publicKey []byte) (ConsensusGroup, error)
	Update(ctx context.Context, id string, data interface{}) (ConsensusGroup, error)
	Create(ctx context.Context, ConsensusGroup ConsensusGroup) (ConsensusGroup, error)
	Delete(ctx context.Context, id string) error
}

// ConsensusGroupRepository represent the ConsensusGroups's repository contract
type ConsensusGroupRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []ConsensusGroup, err error)
	GetByID(ctx context.Context, id string) (ConsensusGroup, error)
	GetByPubKey(ctx context.Context, publicKey []byte) (ConsensusGroup, error)
	Update(ctx context.Context, id string, data interface{}) (ConsensusGroup, error)
	Create(ctx context.Context, ConsensusGroup ConsensusGroup) (ConsensusGroup, error)
	Delete(ctx context.Context, id string) error
}
