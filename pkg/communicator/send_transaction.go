package communicator

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SendTransactionRequest struct {
	Method          string   `json:"method"`
	ContractAddress string   `json:"contract_address"`
	ContractABI     string   `json:"contract_abi"`
	PrivateKeyHex   string   `json:"private_key"` // without "0x" prefix
	Input           []string `json:"input"`       // input parameters for the method
}

type SendTransactionResponse struct {
	TransactionHash string `json:"transaction_hash"`
}

func SendTransaction(ctx context.Context, req SendTransactionRequest) (SendTransactionResponse, error) {
	return sendTransaction(ctx, req)
}

func sendTransaction(ctx context.Context, req SendTransactionRequest) (SendTransactionResponse, error) {
	client, err := ethclient.DialContext(ctx, GetNodeAddress(ctx))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to Ethereum client", slog.Any("err", err))
		return SendTransactionResponse{}, err
	}

	privateKey, err := crypto.HexToECDSA(req.PrivateKeyHex)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to convert private key from hex", slog.Any("err", err))
		return SendTransactionResponse{}, fmt.Errorf("failed to convert private key from hex: %v", err)
	}

	// Derive sender address
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Get the nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get nonce", slog.Any("err", err))
		return SendTransactionResponse{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Gas parameters
	gasLimit := uint64(100000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to suggest gas price", slog.Any("err", err))
		return SendTransactionResponse{}, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// Contract address
	contractAddress := common.HexToAddress(req.ContractAddress)

	callData, _, err := getCallData(ctx, req.ContractABI, req.Method, req.Input)
	if err != nil {
		return SendTransactionResponse{}, fmt.Errorf("failed to get call data: %v", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, callData)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal("Failed to get chain ID:", err)
	}
	if chainID == nil {
		chainID = big.NewInt(1)
	}

	// Sign it
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to sign transaction", slog.Any("err", err))
		return SendTransactionResponse{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send transaction", slog.Any("err", err))
		return SendTransactionResponse{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	return SendTransactionResponse{
		TransactionHash: signedTx.Hash().Hex(),
	}, nil
}
