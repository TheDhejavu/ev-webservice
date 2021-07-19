package entity

import (
	"context"

	blockchain "github.com/thedhejavu/ev-blockchain-protocol/core"
	"github.com/thedhejavu/ev-blockchain-protocol/rpc"
)

type Blockchain struct{}

type BlockchainService interface {
	GetBlockchain(ctx context.Context) ([]map[string]interface{}, error)
}

type BlockchainRepository interface {
	FindTxWithTxOutput(pubkey, ttype string) blockchain.Transaction
	QueryResults(pubkey string) (rpc.Result, error)
	QueryBlockchain() (rpc.Result, error)
	QueryUnUsedBallotTxs(pubkey string) []map[string]blockchain.TxBallotOutput
	GetTransaction(id string) (rpc.Result, error)
	StartElection(pubkey, title, description string, totalPeople int64, candidates [][]byte, groupSigners []string) (rpc.Result, error)
	StopElection(pubkey string, groupSigners []string) (rpc.Result, error)
	StartAccreditation(pubkey, txElectionOutId string, groupSigners []string) (rpc.Result, error)
	StopAccreditation(pubkey, txElectionOutId string, txAcOutId string, groupSigners []string) (rpc.Result, error)
	StartVoting(pubkey string, txElectionOutId string, groupSigners []string) (rpc.Result, error)
	StopVoting(pubkey, txElectionOutId, txVotingOutId string, groupSigners []string) (rpc.Result, error)
	CreateBallot(userId, pubkey, txElectionOutId string, groupSigners []string) (rpc.Result, error)
	CastBallot(userId, pubkey, txElectionOutId, candidatePubkey string) (rpc.Result, error)
}
