package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/PumpkinSeed/letherscan/pkg/communicator"
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
	mux.HandleFunc("GET /transaction/{hash}", getTransactionByHash)

	slog.Info("Starting server on port 8080", slog.String("ethereum_client_url", etherClient.Address))
	http.ListenAndServe(":8080", withCORS(mux))
}

func getTransactionByHash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	respStruct, err := etherClient.GetTransactionByHash(ctx, communicator.GetTransactionByHashRequest{
		Hash: r.PathValue("hash"),
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

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set your CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // or restrict to specific origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the original handler
		h.ServeHTTP(w, r)
	})
}
