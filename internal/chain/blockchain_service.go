package chain

import (
	"context"
	"encoding/json"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type blockchainService struct {
	blockchainRepo entity.BlockchainRepository
	logger         log.Logger
}

func NewBlockchainService(blockchainRepo entity.BlockchainRepository, logger log.Logger) entity.BlockchainService {
	return &blockchainService{
		blockchainRepo: blockchainRepo,
		logger:         logger,
	}
}

func (service *blockchainService) GetBlockchain(ctx context.Context) (res []map[string]interface{}, err error) {
	result, err := service.blockchainRepo.QueryBlockchain()
	// fmt.Println(result)
	inrec, _ := json.Marshal(result.Data)
	json.Unmarshal(inrec, &res)

	return
}
