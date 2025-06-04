package communicator

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGetLatestNBlock(t *testing.T) {
	// Create a mock client
	client := NewClient("http://localhost:8545")

	// Define the context
	ctx := context.Background()

	mes := time.Now()
	// Call the function with a valid client and context
	resp, err := getLatestNBlock(ctx, client, GetLatestNBlockRequest{
		NumberOfBlocks: 7,
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	fmt.Println(resp)

	fmt.Println(time.Since(mes))
}
