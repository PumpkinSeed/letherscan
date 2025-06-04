package communicator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type GetLatestNBlockRequest struct {
	// Number of blocks to retrieve
	NumberOfBlocks int64 `json:"number_of_blocks"`

	// Block number to start from (reversed order)
	BlockNumber int64 `json:"block_number"`
}
type GetLatestNBlockResponse struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Header       *types.Header `json:"header"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Hash                 string `json:"hash"`
	Nonce                string `json:"nonce"`
	BlockHash            string `json:"blockHash"`
	BlockNumber          string `json:"blockNumber"`
	TransactionIndex     string `json:"transactionIndex"`
	From                 string `json:"from"`
	To                   string `json:"to"`
	Value                string `json:"value"`
	GasPrice             string `json:"gasPrice"`
	Gas                  string `json:"gas"`
	Input                string `json:"input"`
	V                    string `json:"v"`
	R                    string `json:"r"`
	S                    string `json:"s"`
	ChainId              string `json:"chainId"`
	Type                 string `json:"type"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
}

func getLatestNBlock(ctx context.Context, client *Client, req GetLatestNBlockRequest) (GetLatestNBlockResponse, error) {
	blockNumber, err := client.eth.BlockNumber(ctx)
	if err != nil {
		return GetLatestNBlockResponse{}, err
	}

	var response GetLatestNBlockResponse

	for i := int64(0); i < req.NumberOfBlocks; i++ {
		block, err := client.eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)-i))
		if err != nil {
			return GetLatestNBlockResponse{}, err
		}

		var transactions []Transaction
		for _, transaction := range block.Transactions() {
			transactions = append(transactions, Transaction{
				Hash: transaction.Hash().Hex(),
			})
		}
		response.Blocks = append(response.Blocks, Block{
			Header:       block.Header(),
			Transactions: transactions,
		})
	}

	return response, nil
}
