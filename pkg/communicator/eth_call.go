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
	//client, err := ethclient.DialContext(GetNodeAddress(ctx))
	client, err := ethclient.DialContext(ctx, "https://eth-mainnet.g.alchemy.com/v2/C7p6saTF6QMxOdXiMOUlIPH-0Sb4PKgC")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to Ethereum client", slog.Any("err", err))
		return ETHCallResponse{}, err
	}

	callData, err := getCallData(ctx, req.ContractABI, req.Method, req.Input)
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

	//decoded := make(map[string]interface{})
	//if err := method.Outputs.UnpackIntoMap(decoded, result); err != nil {
	//	slog.ErrorContext(ctx, "Failed to unpack result", slog.Any("method", method), slog.Any("err", err), slog.Any("result", hex.EncodeToString(result)))
	//	return ETHCallResponse{}, fmt.Errorf("failed to unpack result: %v", err)
	//}
	//outputs, err := method.Outputs.Unpack(result)
	//if err != nil {
	//	slog.ErrorContext(ctx, "Failed to unpack result", slog.Any("method", method), slog.Any("err", err))
	//	return ETHCallResponse{}, fmt.Errorf("failed to unpack result: %v", err)
	//}
	//decoded := make(map[string]interface{})
	//for i, output := range method.Outputs {
	//	switch output.Type.String() {
	//	case "address":
	//		decoded[output.Name] = common.BytesToAddress(outputs[i].([]byte))
	//	case "uint256":
	//		if val, ok := outputs[i].(*big.Int); ok {
	//			decoded[output.Name] = val.String()
	//		} else {
	//			decoded[output.Name] = outputs[i]
	//		}
	//	case "bytes32":
	//		if b32, ok := outputs[i].([32]byte); ok {
	//			decoded[output.Name] = hex.EncodeToString(b32[:])
	//		} else {
	//			decoded[output.Name] = outputs[i]
	//		}
	//	case "string":
	//		if str, ok := outputs[i].(string); ok {
	//			decoded[output.Name] = str
	//		} else {
	//			decoded[output.Name] = outputs[i]
	//		}
	//	}
	//}

	return ETHCallResponse{
		RawResponse: hex.EncodeToString(result),
		//Decoded:     decoded,
	}, nil
}
