package accrediation

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/facialrecognition"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type accreditationService struct {
	blockchainRepo entity.BlockchainRepository
	electionRepo   entity.ElectionRepository
	consensusGroup entity.ConsensusGroupRepository
	logger         log.Logger
}

func NewAccreditationService(
	blockchainRepo entity.BlockchainRepository,
	electionRepo entity.ElectionRepository,
	consensusGroup entity.ConsensusGroupRepository,
	logger log.Logger,
) entity.AccreditationService {
	return &accreditationService{
		blockchainRepo: blockchainRepo,
		electionRepo:   electionRepo,
		consensusGroup: consensusGroup,
		logger:         logger,
	}
}

func (service *accreditationService) CreateBallot(ctx context.Context, electionId, userId, facialImagePath string) (res entity.Accreditation, err error) {
	election, err := service.electionRepo.GetByID(ctx, electionId)
	if err != nil {
		return
	}
	var groupSigners []string
	groups, _ := service.consensusGroup.Fetch(ctx, map[string]string{
		"id": election.Country.ID.Hex(),
	})
	for i := 0; i < len(groups); i++ {
		name := fmt.Sprintf("consensus_%s", groups[i].ID.Hex())
		groupSigners = append(groupSigners, name)
	}

	fg := facialrecognition.NewFacialRecogntion(service.logger)
	result, err := fg.Verify(userId, facialImagePath)
	fmt.Println(result)
	if err != nil {
		return
	}

	if result["result"] != userId {
		err = errors.New("You facial image is not valid, please try again.")
		return
	}
	service.logger.Info("Successful facial recognition")

	key := base64.StdEncoding.EncodeToString([]byte(election.Pubkey))

	bResult, err := service.blockchainRepo.CreateBallot(
		userId,
		key,
		election.TxOutRef,
		groupSigners,
	)
	var respData map[string]interface{}

	inrec, _ := json.Marshal(bResult.Data)
	json.Unmarshal(inrec, &respData)

	fmt.Println(bResult)
	election.Pubkey = key
	res = entity.Accreditation{
		Election:      election,
		BallotTxOutId: fmt.Sprintf("%s", respData["tx_id"]),
	}
	if err != nil {
		return
	}
	return
}

func (service *accreditationService) Start(ctx context.Context, id string) (res entity.Accreditation, err error) {
	return
}
func (service *accreditationService) Stop(ctx context.Context, id string) (res entity.Accreditation, err error) {
	return
}
