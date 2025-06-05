package communicator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type DecodeContractCallDataRequest struct {
	ContractABI string `json:"contract_abi"`
	InputData   string `json:"input_data"`
}

type DecodeContractCallDataResponse struct {
	FunctionName string                 `json:"function_name"`
	Args         map[string]interface{} `json:"args"`
}

func DecodeContractCallData(ctx context.Context, req DecodeContractCallDataRequest) (DecodeContractCallDataResponse, error) {
	return decodeContractCallData(ctx, req)
}

func decodeContractCallData(ctx context.Context, req DecodeContractCallDataRequest) (DecodeContractCallDataResponse, error) {
	data := common.FromHex(req.InputData)

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(req.ContractABI))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse contract ABI", slog.Any("err", err))
		return DecodeContractCallDataResponse{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Extract the function selector (first 4 bytes)
	selector := data[:4]
	payload := data[4:]

	var response = DecodeContractCallDataResponse{
		FunctionName: "",
		Args:         make(map[string]interface{}),
	}
	// Try to match it against all functions
	for name, method := range parsedABI.Methods {
		if string(method.ID) == string(selector) {
			response.FunctionName = name

			// Decode the parameters
			args := make(map[string]interface{})
			if err := method.Inputs.UnpackIntoMap(args, payload); err != nil {
				slog.ErrorContext(ctx, "Failed to decode contract call data", slog.Any("function", name), slog.Any("err", err))
				return DecodeContractCallDataResponse{}, fmt.Errorf("failed to decode args: %v", err)
			}

			for k, v := range args {
				response.Args[k] = v
			}
			return response, nil
		}
	}

	slog.ErrorContext(ctx, "No matching function found for selector", slog.Any("selector", selector))
	return DecodeContractCallDataResponse{}, fmt.Errorf("no matching function found for selector %x", selector)
}
