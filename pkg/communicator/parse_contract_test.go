package communicator

import (
	"context"
	"testing"
)

func TestParseContractABI(t *testing.T) {
	// Call the function to parse the ABI
	parsedABI, err := parseContractABI(context.Background(), ParseContractABIRequest{
		ContractABI: contractABI,
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	t.Log(parsedABI)
}
