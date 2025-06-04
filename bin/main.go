package main

import (
	"encoding/json"
	"github.com/PumpkinSeed/letherscan/pkg/communicator"
	"log/slog"
	"net/http"
	"os"
)

const (
	EthereumClientURLEnv = "ETHEREUM_CLIENT_URL"
)

var etherClient *communicator.Client

func init() {
	var defaultURL = "http://localhost:8545"
	if defaultURLEnv := os.Getenv(EthereumClientURLEnv); defaultURLEnv != "" {
		defaultURL = defaultURLEnv
	}
	// Initialize the Ethereum client
	etherClient = communicator.NewClient(defaultURL)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /blocks", getBlocks)

	slog.Info("Starting server on port 8080", slog.String("ethereum_client_url", etherClient.Address))
	http.ListenAndServe(":8080", mux)
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	respStruct, err := etherClient.GetLatestNBlock(ctx, communicator.GetLatestNBlockRequest{
		NumberOfBlocks: 7,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(respStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
