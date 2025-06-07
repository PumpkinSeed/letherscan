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

func getCallData(ctx context.Context, contractABI, selectedMethod string, inputStr []string) ([]byte, abi.Method, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse contract ABI", slog.Any("err", err))
		return nil, abi.Method{}, fmt.Errorf("failed to parse ABI: %v", err)
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
		return nil, abi.Method{}, fmt.Errorf("failed to pack input: %v", err)
	}
	callData := append(method.ID, data...)

	return callData, method, nil
}

func parseResult(ctx context.Context, method abi.Method, result []byte) (map[string]interface{}, error) {
	outputs, err := method.Outputs.Unpack(result)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to unpack result", slog.Any("method", method), slog.Any("err", err))
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}
	decoded := make(map[string]interface{})
	for i, output := range method.Outputs {
		name := output.Name
		if name == "" {
			name = fmt.Sprintf("output_%d", i)
		}
		switch output.Type.String() {
		case "address":
			decoded[name] = outputs[i].(common.Address)
		case "uint256", "uint48", "uint64", "uint32", "uint16", "uint8":
			if val, ok := outputs[i].(*big.Int); ok {
				decoded[name] = val.String()
			} else {
				decoded[name] = outputs[i]
			}
		case "bytes32":
			if b32, ok := outputs[i].([32]byte); ok {
				decoded[name] = hex.EncodeToString(b32[:])
			} else {
				decoded[name] = outputs[i]
			}
		case "string":
			if str, ok := outputs[i].(string); ok {
				decoded[name] = str
			} else {
				decoded[name] = outputs[i]
			}
		}
	}

	return decoded, nil
}
