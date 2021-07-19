package election

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"github.com/workspace/evoting/ev-webservice/wallet"
)

type electionService struct {
	electionRepo   entity.ElectionRepository
	blockchainRepo entity.BlockchainRepository
	consensusGroup entity.ConsensusGroupRepository
	logger         log.Logger
}

func NewElectionService(
	electionRepo entity.ElectionRepository,
	blockchainRepo entity.BlockchainRepository,
	consensusGroup entity.ConsensusGroupRepository,
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
	value, _ := service.Exists(ctx, map[string]interface{}{"_id": id}, nil)
	if value == false {
		return res, entity.ErrNotFound
	}

	res, err = service.electionRepo.Update(ctx, id, data)
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
	var candidates [][]byte
	var groupSigners []string
	for i := 0; i < len(election.Candidates); i++ {
		candidate := election.Candidates[i]
		w := wallet.MakeWalletGroup()
		candidate.Pubkey = w.Main.PublicKey
		candidates = append(candidates, w.Main.PublicKey)
	}

	groups, _ := service.consensusGroup.Fetch(ctx, map[string]string{
		"id": election.Country.Hex(),
	})
	for i := 0; i < len(groups); i++ {
		name := fmt.Sprintf("consensus_%s", groups[i].ID.Hex())
		groupSigners = append(groupSigners, name)
	}

	pubkey := base64.StdEncoding.EncodeToString(election.Pubkey)
	// fmt.Println(groupSigners, pubkey)
	result, _ := service.blockchainRepo.StartElection(
		pubkey,
		election.Title,
		election.Description,
		100,
		candidates,
		groupSigners,
	)
	var v map[string]interface{}
	inrec, _ := json.Marshal(result.Data)
	json.Unmarshal(inrec, &v)

	fmt.Println(v["tx_id"])
	idStart := v["tx_id"]

	election.TxOutRef = fmt.Sprintf("%s", idStart)

	res, err = service.electionRepo.Create(ctx, *election)
	if err != nil {
		return
	}

	res.Pubkey = string(wallet.Base58Encode([]byte(res.Pubkey)))

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

func (services *electionService) GetResult(ctx context.Context, filter map[string]interface{}) (res []entity.ElectionRead, err error) {
	return
}
