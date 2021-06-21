package entity

import "context"

type ConsensusGroup struct {
	Name      string  `json:"name" bson:"name,omitempty"`
	ServerUrl string  `json:"server_url" bson:"server_url,omitempty"`
	PublicKey []byte  `json:"public_key" bson:"public_key,omitempty"`
	Country   Country `json:"country" bson:"country,omitempty"`
}

//  ConsensusGroupService represent the ConsensusGroups's usecase
type ConsensusGroupService interface {
	Fetch(ctx context.Context, filter interface{}) (res []ConsensusGroup, err error)
	GetByID(ctx context.Context, id string) (ConsensusGroup, error)
	GetByPubKey(ctx context.Context, public_key []byte) (ConsensusGroup, error)
	Update(ctx context.Context, id string, data interface{}) (ConsensusGroup, error)
	Store(ctx context.Context, ConsensusGroup ConsensusGroup) (ConsensusGroup, error)
	Delete(ctx context.Context, id string) error
}

// ConsensusGroupRepository represent the ConsensusGroups's repository contract
type ConsensusGroupRepository interface {
	Fetch(ctx context.Context, filter interface{}) (res []ConsensusGroup, err error)
	GetByID(ctx context.Context, id string) (ConsensusGroup, error)
	tByPubKey(ctx context.Context, public_key []byte) (ConsensusGroup, error)
	Update(ctx context.Context, id string, data interface{}) (ConsensusGroup, error)
	Store(ctx context.Context, ConsensusGroup ConsensusGroup) (ConsensusGroup, error)
	Delete(ctx context.Context, id string) error
}
