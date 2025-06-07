package communicator

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func getCallData(ctx context.Context, contractABI, selectedMethod string, inputStr []string) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse contract ABI", slog.Any("err", err))
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	var method abi.Method
	for _, m := range parsedABI.Methods {
		if strings.Contains(m.Name, selectedMethod) {
			method = m
		}
	}

	var convertedArgs []interface{}
	for i, input := range method.Inputs {
		switch input.Type.String() {
		case "address":
			convertedArgs = append(convertedArgs, common.HexToAddress(inputStr[i]))
		case "uint256":
			val := new(big.Int)
			val.SetString(inputStr[i], 10)
			convertedArgs = append(convertedArgs, val)
		case "bytes32":
			var b32 [32]byte
			copy(b32[:], inputStr[i])
			convertedArgs = append(convertedArgs, b32)
		case "string":
			convertedArgs = append(convertedArgs, inputStr[i])
		default:
			slog.ErrorContext(ctx, "Unsupported input type", slog.Any("type", input.Type.String()))
		}
	}

	data, err := method.Inputs.Pack(convertedArgs...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to pack input", slog.Any("err", err))
		return nil, fmt.Errorf("failed to pack input: %v", err)
	}
	callData := append(method.ID, data...)
	slog.InfoContext(ctx, "Generated call data", slog.Any("method", method.Name))
	slog.InfoContext(ctx, "Generated call data", slog.Any("data", hex.EncodeToString(callData)))
	slog.InfoContext(ctx, "Generated call data", slog.Any("method_id", hex.EncodeToString(method.ID)))
	return callData, nil
}
