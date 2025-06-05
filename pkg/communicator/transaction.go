package communicator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

func GetTransactionByHash(ctx context.Context, req GetTransactionByHashRequest) (Transaction, error) {
	return getTransactionByHash(ctx, req)
}

func getTransactionByHash(ctx context.Context, req GetTransactionByHashRequest) (Transaction, error) {
	client, err := ethclient.Dial(GetNodeAddress(ctx))
	if err != nil {
		return Transaction{}, err
	}

	transaction, isPending, err := client.TransactionByHash(ctx, common.HexToHash(req.Hash))
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
	chainID := transaction.ChainId()
	if chainID == nil || chainID.Int64() == 0 {
		chainID = big.NewInt(1)
	}
	signer := types.LatestSignerForChainID(chainID)
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
		To:               safeHexAddress(transaction.To()),
		Value:            safeBigIntToString(transaction.Value()),
		GasPrice:         safeBigIntToString(transaction.GasPrice()),
		Gas:              transaction.Gas(),
		Input:            fmt.Sprintf("0x%x", transaction.Data()),
		V:                safeBigIntToString(v),
		S:                safeBigIntToString(s),
		R:                safeBigIntToString(r),
		ChainId:          safeBigIntToString(chainID),
		Type:             parseTransactionType(transaction.Type()),
		Method:           parseMethod(transaction),
		IsPending:        false,
	}, nil
}

func safeHexAddress(addr *common.Address) string {
	if addr == nil {
		return ""
	}
	return addr.Hex()
}

func safeBigIntToString(value *big.Int) string {
	if value == nil {
		return "0"
	}
	return value.String()
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
