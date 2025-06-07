package communicator

import (
	"context"
	"testing"
)

func TestETHCall(t *testing.T) {
	// Define the context
	ctx := context.Background()

	// Call the function with a valid client and context
	resp, err := ethCall(ctx, ETHCallRequest{
		Method:          "balanceOf",
		ContractAddress: "0x514910771AF9Ca656af840dff83E8264EcF986CA", // Chainlink token contract address
		ContractABI:     contractABI,
		Input:           []string{"0x9491A3757A98e53BE0d1c14834a6e2Da0B4Dc527"}, // random address for testing
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Response: %v", resp)
}
