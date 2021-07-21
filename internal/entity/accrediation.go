package entity

import (
	"context"
)

//Accreditation represent the Accreditation entity
type Accreditation struct {
	Election      ElectionRead `json:"election" `
	BallotTxOutId string       `json:"ballot_txout_id"`
}

//  AccreditationService represent the Accreditations's usecase
type AccreditationService interface {
	CreateBallot(ctx context.Context, electionId, userId, facialImagePath string) (Accreditation, error)
	Start(ctx context.Context, electionId string) (Accreditation, error)
	Stop(ctx context.Context, electionId string) (Accreditation, error)
}
