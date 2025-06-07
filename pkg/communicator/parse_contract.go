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
	Methods []Method `json:"methods"`
}

type Method struct {
	Name            string   `json:"name"`
	StateMutability string   `json:"state_mutability"`
	Inputs          []string `json:"inputs"`
	Outputs         []string `json:"outputs"`
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
			var inputs []string
			for _, input := range method.Inputs {
				inputs = append(inputs, fmt.Sprintf("%s %s", input.Type.String(), input.Name))
			}
			var outputs []string
			for _, output := range method.Outputs {
				outputs = append(outputs, fmt.Sprintf("%s %s", output.Type.String(), output.Name))
			}
			response.Methods = append(response.Methods, Method{
				Name:            method.Name,
				StateMutability: method.StateMutability,
				Inputs:          inputs,
				Outputs:         outputs,
			})
		}
	}

	return response, nil
}
