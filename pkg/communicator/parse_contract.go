package communicator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ParseContractABIRequest struct {
	ContractABI           string `json:"contract_abi"`
	StateMutabilityFilter string `json:"state_mutability_filter"`
}

type ParseContractABIResponse struct {
	Methods []string `json:"methods"`
}

func ParseContractABI(ctx context.Context, req ParseContractABIRequest) (ParseContractABIResponse, error) {
	return parseContractABI(ctx, req)
}

func parseContractABI(ctx context.Context, req ParseContractABIRequest) (ParseContractABIResponse, error) {
	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(req.ContractABI))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse contract ABI", slog.Any("err", err))
		return ParseContractABIResponse{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	var response ParseContractABIResponse
	for _, method := range parsedABI.Methods {
		if method.StateMutability == req.StateMutabilityFilter || req.StateMutabilityFilter == "" {
			response.Methods = append(response.Methods, strings.ReplaceAll(method.String(), "function ", ""))
		}
	}

	return response, nil
}
