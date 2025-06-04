package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/PumpkinSeed/letherscan/pkg/communicator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed build/*
var embeddedFiles embed.FS

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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	distFS, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		log.Fatal(err)
	}

	r.Handle("/*", http.FileServer(http.FS(distFS)))

	r.Get("/blocks", getBlocks)
	r.Get("/transaction/{hash}", getTransactionByHash)

	slog.Info("starting server", slog.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func getTransactionByHash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	respStruct, err := etherClient.GetTransactionByHash(ctx, communicator.GetTransactionByHashRequest{
		Hash: chi.URLParam(r, "hash"),
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
