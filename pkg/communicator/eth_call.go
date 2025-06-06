package communicator

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	//client, err := ethclient.DialContext(GetNodeAddress(ctx))
	client, err := ethclient.DialContext(ctx, "https://eth-mainnet.g.alchemy.com/v2/C7p6saTF6QMxOdXiMOUlIPH-0Sb4PKgC")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to Ethereum client", slog.Any("err", err))
		return ETHCallResponse{}, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(req.ContractABI))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse contract ABI", slog.Any("err", err))
		return ETHCallResponse{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	var method abi.Method
	for _, m := range parsedABI.Methods {
		if strings.Contains(m.String(), req.Method) {
			method = m
		}
	}

	var convertedArgs []interface{}
	for i, input := range method.Inputs {
		switch input.Type.String() {
		case "address":
			convertedArgs = append(convertedArgs, common.HexToAddress(req.Input[i]))
		case "uint256":
			val := new(big.Int)
			val.SetString(req.Input[i], 10)
			convertedArgs = append(convertedArgs, val)
		case "bytes32":
			var b32 [32]byte
			copy(b32[:], req.Input[i])
			convertedArgs = append(convertedArgs, b32)
		case "string":
			convertedArgs = append(convertedArgs, req.Input[i])
		default:
			slog.ErrorContext(ctx, "Unsupported input type", slog.Any("type", input.Type.String()))
		}
	}

	data, err := method.Inputs.Pack(convertedArgs...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to pack input", slog.Any("err", err))
		return ETHCallResponse{}, fmt.Errorf("failed to pack input: %v", err)
	}
	callData := append(method.ID, data...)

	contractAddress := common.HexToAddress(req.ContractAddress)
	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to call contract", slog.Any("method", method), slog.Any("err", err))
		return ETHCallResponse{}, err
	}
	fmt.Println(hex.EncodeToString(result))

	outputs, err := method.Outputs.Unpack(result)
	if err != nil {
		log.Fatal("failed to unpack result:", err)
	}
	decoded := make(map[string]interface{})
	for i, output := range method.Outputs {
		switch output.Type.String() {
		case "address":
			decoded[output.Name] = common.BytesToAddress(outputs[i].([]byte))
		case "uint256":
			if val, ok := outputs[i].(*big.Int); ok {
				decoded[output.Name] = val.String()
			} else {
				decoded[output.Name] = outputs[i]
			}
		case "bytes32":
			if b32, ok := outputs[i].([32]byte); ok {
				decoded[output.Name] = hex.EncodeToString(b32[:])
			} else {
				decoded[output.Name] = outputs[i]
			}
		case "string":
			if str, ok := outputs[i].(string); ok {
				decoded[output.Name] = str
			} else {
				decoded[output.Name] = outputs[i]
			}
		}
	}

	return ETHCallResponse{
		RawResponse: hex.EncodeToString(result),
		Decoded:     decoded,
	}, nil
}
