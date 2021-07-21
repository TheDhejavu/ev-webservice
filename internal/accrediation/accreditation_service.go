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
	"github.com/workspace/evoting/ev-webservice/wallet"
)

type accreditationService struct {
	blockchainService        entity.BlockchainService
	electionRepo          entity.ElectionRepository
	consensusGroupService entity.ConsensusGroupService
	logger                log.Logger
}

func NewAccreditationService(
	blockchainService entity.BlockchainService,
	electionRepo entity.ElectionRepository,
	consensusGroupService entity.ConsensusGroupService,
	logger log.Logger,
) entity.AccreditationService {
	return &accreditationService{
		blockchainService:        blockchainService,
		electionRepo:          electionRepo,
		consensusGroupService: consensusGroupService,
		logger:                logger,
	}
}

func (service *accreditationService) CreateBallot(ctx context.Context, electionId, userId, facialImagePath string) (res entity.Accreditation, err error) {
	election, err := service.electionRepo.GetByID(ctx, electionId)
	if err != nil {
		return
	}

	groupSigners, err := service.consensusGroupService.GetIDs(ctx, electionId)
	if err != nil {
		return
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

	bResult, err := service.blockchainService.CreateBallot(
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

	return
}

func (service *accreditationService) Start(ctx context.Context, electionId string) (res entity.Accreditation, err error) {
	election, err := service.electionRepo.Get(ctx, map[string]interface{}{
		"id": electionId,
	})
	if err != nil {
		return
	}

	groupSigners, err := service.consensusGroupService.GetIDs(ctx, election.Country.Hex())
	if err != nil {
		return
	}
	pubkey := base64.StdEncoding.EncodeToString([]byte(election.Pubkey))

	result, err := service.blockchainService.StartAccreditation(
		pubkey,
		election.TxOutRef,
		groupSigners,
	)
	if err != nil {
		return
	}
	var v map[string]interface{}
	inrec, _ := json.Marshal(result.Data)
	json.Unmarshal(inrec, &v)

	fmt.Println(v["tx_id"])
	txId := fmt.Sprintf("%s", v["tx_id"])
	election.AccreditationAt.TxStartRef = txId
	election.Phase = "accreditation"

	updatedElection, err := service.electionRepo.Update(ctx, electionId, election)
	if err != nil {
		return
	}
	updatedElection.Pubkey = string(wallet.Base58Encode([]byte(election.Pubkey)))

	res = entity.Accreditation{
		Election: updatedElection,
	}
	return
}
func (service *accreditationService) Stop(ctx context.Context, electionId string) (res entity.Accreditation, err error) {
	election, err := service.electionRepo.Get(ctx, map[string]interface{}{
		"id": electionId,
	})
	if err != nil {
		return
	}

	groupSigners, err := service.consensusGroupService.GetIDs(ctx, election.Country.Hex())
	if err != nil {
		return
	}

	pubkey := base64.StdEncoding.EncodeToString([]byte(election.Pubkey))

	fmt.Println("REF_OUT", election.AccreditationAt.TxStartRef)
	result, err := service.blockchainService.StopAccreditation(
		pubkey,
		election.TxOutRef,
		election.AccreditationAt.TxStartRef,
		groupSigners,
	)
	if err != nil {
		return
	}
	var v map[string]interface{}
	inrec, _ := json.Marshal(result.Data)
	json.Unmarshal(inrec, &v)

	fmt.Println(v["tx_id"])
	txId := fmt.Sprintf("%s", v["tx_id"])

	election.AccreditationAt.TxEndRef = txId
	election.Phase = "accreditation_end"

	updatedElection, err := service.electionRepo.Update(ctx, electionId, election)
	if err != nil {
		return
	}
	updatedElection.Pubkey = string(wallet.Base58Encode([]byte(election.Pubkey)))

	res = entity.Accreditation{
		Election: updatedElection,
	}
	return
}
