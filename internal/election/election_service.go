package election

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"github.com/workspace/evoting/ev-webservice/wallet"
)

type electionService struct {
	electionRepo   entity.ElectionRepository
	blockchainRepo entity.BlockchainService
	consensusGroup entity.ConsensusGroupService
	logger         log.Logger
}

func NewElectionService(
	electionRepo entity.ElectionRepository,
	blockchainRepo entity.BlockchainService,
	consensusGroup entity.ConsensusGroupService,
	logger log.Logger,
) entity.ElectionService {
	return &electionService{
		electionRepo:   electionRepo,
		blockchainRepo: blockchainRepo,
		consensusGroup: consensusGroup,
		logger:         logger,
	}
}

func (service *electionService) Fetch(ctx context.Context, filter map[string]interface{}) (res []*entity.ElectionRead, err error) {
	res, err = service.electionRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}
	for i := 0; i < len(res); i++ {
		data := res[i]
		data.Pubkey = string(wallet.Base58Encode([]byte(data.Pubkey)))
	}
	return
}
func (service *electionService) GetByID(ctx context.Context, id string) (res entity.ElectionRead, err error) {
	res, err = service.electionRepo.GetByID(ctx, id)
	res.Pubkey = string(wallet.Base58Encode([]byte(res.Pubkey)))
	if err != nil {
		return
	}
	return
}
func (service *electionService) Update(ctx context.Context, id string, data map[string]interface{}) (res entity.ElectionRead, err error) {
	election, err := service.electionRepo.Get(ctx, map[string]interface{}{"_id": id})

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return res, err
		default:
			return res, err
		}
	}

	res, err = service.electionRepo.Update(ctx, id, election)
	if err != nil {
		return
	}
	return
}

func (service *electionService) Create(ctx context.Context, data map[string]interface{}) (res entity.ElectionRead, err error) {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		return
	}

	election := &entity.Election{}
	if err = json.Unmarshal(jsonbody, &election); err != nil {
		return
	}
	// Initialize system identity wallet
	w := wallet.MakeWalletGroup()

	election.Pubkey = w.Main.PublicKey
	for i := 0; i < len(election.Candidates); i++ {
		candidate := election.Candidates[i]
		w := wallet.MakeWalletGroup()
		candidate.Pubkey = w.Main.PublicKey
	}

	res, err = service.electionRepo.Create(ctx, *election)
	if err != nil {
		return
	}

	res.Pubkey = string(wallet.Base58Encode([]byte(res.Pubkey)))

	return
}

func (service *electionService) Start(ctx context.Context, electionId string) (res entity.ElectionRead, err error) {
	election, err := service.electionRepo.Get(ctx, map[string]interface{}{"_id": electionId})

	if err != nil {
		return
	}

	groups, _ := service.consensusGroup.Fetch(ctx, map[string]string{
		"id": election.Country.Hex(),
	})

	if len(groups) == 0 {
		err = errors.New("No consensus group found")
		return
	}
	var candidates [][]byte
	var groupSigners []string
	for i := 0; i < len(groups); i++ {
		name := fmt.Sprintf("consensus_%s", groups[i].ID.Hex())
		groupSigners = append(groupSigners, name)
	}

	for i := 0; i < len(election.Candidates); i++ {
		pubkey := election.Candidates[i].Pubkey
		candidates = append(candidates, pubkey)
	}

	pubkey := base64.StdEncoding.EncodeToString([]byte(election.Pubkey))

	result, err := service.blockchainRepo.StartElection(
		pubkey,
		election.Title,
		election.Description,
		100,
		candidates,
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
	election.TxOutRef = txId

	updatedElection, err := service.electionRepo.Update(ctx, electionId, election)
	if err != nil {
		return
	}

	updatedElection.Pubkey = string(wallet.Base58Encode([]byte(election.Pubkey)))

	res = updatedElection
	return
}

func (service *electionService) Stop(ctx context.Context, electionId string) (election entity.ElectionRead, err error) {
	election, err = service.electionRepo.GetByID(ctx, electionId)
	if err != nil {
		return
	}

	groupSigners, err := service.consensusGroup.GetIDs(ctx, election.Country.ID.Hex())
	if err != nil {
		return
	}
	var candidates [][]byte
	for i := 0; i < len(election.Candidates); i++ {
		pubkey := election.Candidates[i].Pubkey
		candidates = append(candidates, pubkey)
	}

	pubkey := base64.StdEncoding.EncodeToString([]byte(election.Pubkey))

	result, err := service.blockchainRepo.StopElection(
		pubkey,
		groupSigners,
	)
	if err != nil {
		return
	}
	var v map[string]interface{}
	inrec, _ := json.Marshal(result.Data)
	json.Unmarshal(inrec, &v)

	fmt.Println(v["tx_id"])

	election.Pubkey = string(wallet.Base58Encode([]byte(election.Pubkey)))

	return
}
func (service *electionService) Delete(ctx context.Context, id string) (err error) {
	value, _ := service.Exists(ctx, map[string]interface{}{"_id": id}, nil)
	if value == false {
		return entity.ErrNotFound
	}

	err = service.electionRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}

func (service *electionService) Exists(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (res bool, err error) {
	if exclude == nil {
		_, err = service.electionRepo.Get(ctx, filter)
	} else {
		_, err = service.electionRepo.GetWithExclude(ctx, filter, exclude)
	}

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, err
}

func (service *electionService) GetResults(ctx context.Context, id string) (res []*entity.CandidateRead, err error) {
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

	fmt.Println("RESPONSE", results)

	return
}
