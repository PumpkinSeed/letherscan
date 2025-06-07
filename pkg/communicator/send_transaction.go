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

type SendTransactionResponse struct{}

func SendTransaction(ctx context.Context, req SendTransactionRequest) (SendTransactionResponse, error) {
	//client, err := ethclient.DialContext(GetNodeAddress(ctx))
	client, err := ethclient.DialContext(ctx, "https://eth-mainnet.g.alchemy.com/v2/C7p6saTF6QMxOdXiMOUlIPH-0Sb4PKgC")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to Ethereum client", slog.Any("err", err))
		return SendTransactionResponse{}, err
	}

	privateKey, err := crypto.HexToECDSA(req.PrivateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// Derive sender address
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Gas parameters
	gasLimit := uint64(100000)                                    // enough for simple call
	gasPrice, err := client.SuggestGasPrice(context.Background()) // or hardcode if needed
	if err != nil {
		log.Fatal(err)
	}

	// Contract address
	contractAddress := common.HexToAddress(req.ContractAddress)

	callData, _, err := getCallData(ctx, req.ContractABI, req.Method, req.Input)
	if err != nil {
		return SendTransactionResponse{}, fmt.Errorf("failed to get call data: %v", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, callData)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("Failed to get chain ID:", err)
	}
	if chainID == nil {
		chainID = big.NewInt(1) // Default to mainnet if chain ID is not available
	}

	// Sign it
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sent tx: %s\n", signedTx.Hash().Hex())
	return SendTransactionResponse{}, nil
}
