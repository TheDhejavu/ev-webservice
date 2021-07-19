package voting

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type votingService struct {
	blockchainRepo entity.BlockchainRepository
	electionRepo   entity.ElectionRepository
	logger         log.Logger
}

func NewVotingService(
	blockchainRepo entity.BlockchainRepository,
	electionRepo entity.ElectionRepository,
	logger log.Logger) entity.VotingService {
	return &votingService{
		blockchainRepo: blockchainRepo,
		electionRepo:   electionRepo,
		logger:         logger,
	}
}

func (service *votingService) GetResults(ctx context.Context, id string) (res []*entity.CandidateRead, err error) {
	election, err := service.electionRepo.GetByID(ctx, id)

	if err != nil {
		return
	}
	key := base64.StdEncoding.EncodeToString([]byte(election.Pubkey))
	results, err := service.blockchainRepo.QueryResults(key)
	if err != nil {
		return
	}

	var resp map[string]int64

	inrec, _ := json.Marshal(results.Data)
	json.Unmarshal(inrec, &resp)

	res = election.Candidates
	for i := 0; i < len(res); i++ {
		hexStr := hex.EncodeToString(res[i].Pubkey)
		if v, ok := resp[hexStr]; ok {
			res[i].Result = v
		}
	}

	fmt.Println("RESPONSE", res)

	return
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
func (service *votingService) Start(ctx context.Context, id string) (res entity.Vote, err error) {
	return
}
func (service *votingService) Stop(ctx context.Context, id string) (res entity.Vote, err error) {
	return
}
