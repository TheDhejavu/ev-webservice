package entity

import (
	"context"
	"crypto/ecdsa"
	"time"
)

//Accreditation represent the Accreditation entity
type Accreditation struct {
	ID             []byte    `json:"id" bson:"_id,omitempty"`
	PoliticalParty string    `json:"political_party" bson:"political_party,omitempty"`
	Election       Election  `json:"election" bson:"election,omitempty"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at,omitempty"`
}

//  AccreditationService represent the Accreditations's usecase
type AccreditationService interface {
	Fetch(ctx context.Context, filter interface{}) (res []Accreditation, err error)
	GetByID(ctx context.Context, id []byte) (Accreditation, error)
	Execute(ctx context.Context, accreditation Accreditation) (Accreditation, error)
	FindBlockTx(ctx context.Context, privkey ecdsa.PrivateKey) (Accreditation, error)
	Start(ctx context.Context, id string) (Election, error)
	Stop(ctx context.Context, id string) (Election, error)
}
