package communicator

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ETHCallRequest struct {
	Method          string   `json:"method"`
	ContractAddress string   `json:"contract_address"`
	ContractABI     string   `json:"contract_abi"`
	Input           []string `json:"input"`
}

type ETHCallResponse struct {
	RawResponse string                 `json:"raw_response"`
	Decoded     map[string]interface{} `json:"decoded"`
}

func ETHCall(ctx context.Context, req ETHCallRequest) (ETHCallResponse, error) {
	return ethCall(ctx, req)
}

func ethCall(ctx context.Context, req ETHCallRequest) (ETHCallResponse, error) {
	client, err := ethclient.DialContext(ctx, GetNodeAddress(ctx))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to Ethereum client", slog.Any("err", err))
		return ETHCallResponse{}, err
	}

	callData, method, err := getCallData(ctx, req.ContractABI, req.Method, req.Input)
	if err != nil {
		return ETHCallResponse{}, fmt.Errorf("failed to get call data: %v", err)
	}

	slog.InfoContext(ctx, "Calling contract", slog.Any("contract_address", req.ContractAddress))
	contractAddress := common.HexToAddress(req.ContractAddress)
	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to call contract", slog.Any("method", req.Method), slog.Any("err", err))
		return ETHCallResponse{}, err
	}
	if len(result) == 0 {
		slog.ErrorContext(ctx, "No result returned from contract call", slog.Any("method", req.Method))
		return ETHCallResponse{}, fmt.Errorf("no result returned from contract call")
	}

	decoded, err := parseResult(ctx, method, result)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse result", slog.Any("method", req.Method), slog.Any("err", err))
		return ETHCallResponse{}, fmt.Errorf("failed to parse result: %v", err)
	}

	return ETHCallResponse{
		RawResponse: hex.EncodeToString(result),
		Decoded:     decoded,
	}, nil
}
