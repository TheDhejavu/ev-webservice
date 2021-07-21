package voting

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"github.com/workspace/evoting/ev-webservice/wallet"
)

type votingService struct {
	blockchainRepo        entity.BlockchainService
	electionRepo          entity.ElectionRepository
	consensusGroupService entity.ConsensusGroupService
	logger                log.Logger
}

func NewVotingService(
	blockchainRepo entity.BlockchainService,
	electionRepo entity.ElectionRepository,
	consensusGroupService entity.ConsensusGroupService,
	logger log.Logger) entity.VotingService {
	return &votingService{
		blockchainRepo:        blockchainRepo,
		electionRepo:          electionRepo,
		consensusGroupService: consensusGroupService,
		logger:                logger,
	}
}

func (service *votingService) CastVote(ctx context.Context, userId, electionId, candidate string) (res *entity.CandidateRead, err error) {
	election, err := service.electionRepo.GetByID(ctx, electionId)
	if err != nil {
		return
	}

	key := base64.StdEncoding.EncodeToString([]byte(election.Pubkey))

	results, err := service.blockchainRepo.CastBallot(
		userId,
		key,
		election.TxOutRef,
		candidate,
	)
	if err != nil {
		return
	}

	var resp map[string]interface{}

	inrec, _ := json.Marshal(results.Data)
	json.Unmarshal(inrec, &resp)

	fmt.Println("RESPONSE", results)

	for i := 0; i < len(election.Candidates); i++ {
		key := base64.StdEncoding.EncodeToString(election.Candidates[i].Pubkey)
		fmt.Println(key)
		if candidate == key {
			res = election.Candidates[i]
			break
		}
	}
	return
}
func (service *votingService) Start(ctx context.Context, electionId string) (res entity.Vote, err error) {
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

	fmt.Println("REF_OUT", election.VoteAt.TxStartRef)
	result, err := service.blockchainRepo.StartVoting(
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

	election.VoteAt.TxStartRef = txId
	election.Phase = "voting"

	updatedElection, err := service.electionRepo.Update(ctx, electionId, election)
	if err != nil {
		return
	}
	updatedElection.Pubkey = string(wallet.Base58Encode([]byte(election.Pubkey)))

	res = entity.Vote{
		Election: updatedElection,
	}
	return
}
func (service *votingService) Stop(ctx context.Context, electionId string) (res entity.Vote, err error) {
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

	fmt.Println("REF_OUT", election.VoteAt.TxStartRef)
	result, err := service.blockchainRepo.StopVoting(
		pubkey,
		election.TxOutRef,
		election.VoteAt.TxStartRef,
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

	election.VoteAt.TxEndRef = txId
	election.Phase = "voting_end"

	updatedElection, err := service.electionRepo.Update(ctx, electionId, election)
	if err != nil {
		return
	}
	updatedElection.Pubkey = string(wallet.Base58Encode([]byte(election.Pubkey)))

	res = entity.Vote{
		Election: updatedElection,
	}
	return
}
