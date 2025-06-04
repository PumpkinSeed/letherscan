package communicator

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
)

type GetTransactionByHashRequest struct {
	Hash string `json:"hash"`
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

func (client *Client) GetTransactionByHash(ctx context.Context, req GetTransactionByHashRequest) (Transaction, error) {
	return getTransactionByHash(ctx, client, req)
}

func getTransactionByHash(ctx context.Context, client *Client, req GetTransactionByHashRequest) (Transaction, error) {
	transaction, isPending, err := client.eth.TransactionByHash(ctx, common.HexToHash(req.Hash))
	if err != nil {
		return Transaction{}, err
	}
	parsedTransaction, err := parseTransaction(transaction, "", 0)
	if err != nil {
		return Transaction{}, err
	}
	parsedTransaction.IsPending = isPending
	return parsedTransaction, nil
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
