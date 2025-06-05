package communicator

import (
	"context"
)

type ContextKey string

const (
	NodeAddressContextKey ContextKey = "node_address"

	DefaultNodeAddress = "http://localhost:8545" // Default Ethereum node address
)

func SetNodeAddress(ctx context.Context, address string) context.Context {
	return context.WithValue(ctx, NodeAddressContextKey, address)
}

func GetNodeAddress(ctx context.Context) string {
	if address, ok := ctx.Value(NodeAddressContextKey).(string); ok {
		return address
	}
	return DefaultNodeAddress
}
