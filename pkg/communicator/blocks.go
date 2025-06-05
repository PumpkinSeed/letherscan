package communicator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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
	Header       Header        `json:"header"`
	Transactions []Transaction `json:"transactions"`
}

type Header struct {
	ParentHash       string           `json:"parent_hash"`
	UncleHash        string           `json:"uncle_hash"`
	Coinbase         string           `json:"miner"`
	Root             string           `json:"root"`
	TxHash           string           `json:"tx_hash"`
	ReceiptHash      string           `json:"receipt_hash"`
	Bloom            types.Bloom      `json:"bloom"`
	Difficulty       string           `json:"difficulty"`
	Number           string           `json:"number"`
	GasLimit         uint64           `json:"gas_limit"`
	GasUsed          uint64           `json:"gas_used"`
	Time             uint64           `json:"timestamp"`
	Extra            []byte           `json:"extra"`
	MixDigest        string           `json:"mix_digest"`
	Nonce            types.BlockNonce `json:"nonce"`
	BaseFee          string           `json:"base_fee"`
	WithdrawalsHash  string           `json:"withdrawals_hash"`
	BlobGasUsed      uint64           `json:"blob_gas_used"`
	ExcessBlobGas    uint64           `json:"excess_blob_gas"`
	ParentBeaconRoot string           `json:"parent_beacon_root"`
	RequestsHash     string           `json:"requests_hash"`
}

func GetLatestNBlock(ctx context.Context, req GetLatestNBlockRequest) (GetLatestNBlockResponse, error) {
	return getLatestNBlock(ctx, req)
}

func getLatestNBlock(ctx context.Context, req GetLatestNBlockRequest) (GetLatestNBlockResponse, error) {
	client, err := ethclient.Dial(GetNodeAddress(ctx))
	if err != nil {
		return GetLatestNBlockResponse{}, err
	}

	if req.BlockNumber == 0 {
		blockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			return GetLatestNBlockResponse{}, err
		}
		req.BlockNumber = int64(blockNumber)
	}

	var response GetLatestNBlockResponse

	for i := int64(0); i < req.NumberOfBlocks; i++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(req.BlockNumber-i))
		if err != nil {
			return GetLatestNBlockResponse{}, err
		}

		var transactions []Transaction
		for j, transaction := range block.Transactions() {
			parsedTransaction, err := parseTransaction(transaction, block.Number().String(), int64(j))
			if err != nil {
				return GetLatestNBlockResponse{}, err
			}
			transactions = append(transactions, parsedTransaction)
		}
		response.Blocks = append(response.Blocks, Block{
			Header:       parseHeader(block.Header()),
			Transactions: transactions,
		})
	}

	return response, nil
}

func parseHeader(header *types.Header) Header {
	blobGasUsed := uint64(0)
	if header.BlobGasUsed != nil {
		blobGasUsed = *header.BlobGasUsed
	}

	excessBlobGas := uint64(0)
	if header.ExcessBlobGas != nil {
		excessBlobGas = *header.ExcessBlobGas
	}
	return Header{
		ParentHash:       header.ParentHash.Hex(),
		UncleHash:        header.UncleHash.Hex(),
		Coinbase:         header.Coinbase.Hex(),
		Root:             header.Root.Hex(),
		TxHash:           header.TxHash.Hex(),
		ReceiptHash:      header.ReceiptHash.Hex(),
		Bloom:            header.Bloom,
		Difficulty:       header.Difficulty.String(),
		Number:           header.Number.String(),
		GasLimit:         header.GasLimit,
		GasUsed:          header.GasUsed,
		Time:             header.Time,
		Extra:            header.Extra,
		MixDigest:        header.MixDigest.Hex(),
		Nonce:            header.Nonce,
		BaseFee:          header.BaseFee.String(),
		WithdrawalsHash:  header.WithdrawalsHash.Hex(),
		BlobGasUsed:      blobGasUsed,
		ExcessBlobGas:    excessBlobGas,
		ParentBeaconRoot: header.ParentBeaconRoot.Hex(),
		RequestsHash:     header.ReceiptHash.Hex(),
	}
}

func blockHashToString(blockHash []common.Hash) []string {
	var hashes []string
	for _, hash := range blockHash {
		hashes = append(hashes, hash.Hex())
	}
	return hashes
}
