package communicator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

type Transaction struct {
	Hash             string   `json:"hash"`
	Nonce            uint64   `json:"nonce"`
	BlockHash        []string `json:"block_hash"`
	BlockNumber      string   `json:"block_number"`
	TransactionIndex int64    `json:"transaction_index"`
	From             string   `json:"from"`
	To               string   `json:"to"`
	Value            string   `json:"value"`
	GasPrice         string   `json:"gas_price"`
	Gas              uint64   `json:"gas"`
	Input            string   `json:"input"`
	V                string   `json:"v"`
	R                string   `json:"r"`
	S                string   `json:"s"`
	ChainId          string   `json:"chain_id"`
	Type             string   `json:"type"`
	Method           string   `json:"method"`

	IsPending bool `json:"isPending"`
}

func (c *Client) GetLatestNBlock(ctx context.Context, req GetLatestNBlockRequest) (GetLatestNBlockResponse, error) {
	return getLatestNBlock(ctx, c, req)
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

func parseTransaction(transaction *types.Transaction, blockNumber string, index int64) (Transaction, error) {
	signer := types.LatestSignerForChainID(transaction.ChainId())
	sender, err := types.Sender(signer, transaction)
	if err != nil {
		return Transaction{}, err
	}

	v, r, s := transaction.RawSignatureValues()
	return Transaction{
		Hash:             transaction.Hash().Hex(),
		Nonce:            transaction.Nonce(),
		BlockHash:        blockHashToString(transaction.BlobHashes()),
		BlockNumber:      blockNumber,
		TransactionIndex: index,
		From:             sender.Hex(),
		To:               transaction.To().Hex(),
		Value:            transaction.Value().String(),
		GasPrice:         transaction.GasPrice().String(),
		Gas:              transaction.Gas(),
		Input:            fmt.Sprintf("%x", transaction.Data()),
		V:                v.String(),
		S:                s.String(),
		R:                r.String(),
		ChainId:          transaction.ChainId().String(),
		Type:             parseTransactionType(transaction.Type()),
		Method:           parseMethod(transaction),
		IsPending:        false,
	}, nil
}

func parseTransactionType(t uint8) string {
	switch t {
	case types.LegacyTxType:
		return "legacy"
	case types.AccessListTxType:
		return "access_list"
	case types.DynamicFeeTxType:
		return "dynamic_fee"
	case types.BlobTxType:
		return "blob"
	case types.SetCodeTxType:
		return "set_code"
	}
	return "unknown"
}

func parseMethod(t *types.Transaction) string {
	if t.To() == nil {
		return "contract_creation"
	} else if len(t.Data()) == 0 {
		return "native_transfer"
	} else {
		return "contract_call"
	}
	return "unknown"
}

func blockHashToString(blockHash []common.Hash) []string {
	var hashes []string
	for _, hash := range blockHash {
		hashes = append(hashes, hash.Hex())
	}
	return hashes
}
