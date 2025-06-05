package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/PumpkinSeed/letherscan/pkg/communicator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

//go:embed build/*
var embeddedFiles embed.FS

const (
	NodeAddressHeaderKey = "X-Node-Address"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	r := chi.NewRouter()

	// CORS middleware setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", NodeAddressHeaderKey},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Get Node Address from header and set it in context
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			nodeAddress := r.Header.Get(NodeAddressHeaderKey)
			if nodeAddress != "" {
				r = r.WithContext(communicator.SetNodeAddress(r.Context(), nodeAddress))
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	// Logging middleware
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			slog.InfoContext(
				r.Context(),
				"Request received",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
			)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	distFS, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		log.Fatal(err)
	}

	r.Handle("/*", http.FileServer(http.FS(distFS)))

	r.Get("/blocks", getBlocks)
	r.Get("/transaction/{hash}", getTransactionByHash)
	r.Post("/decode-contract-call-data", decodeContractCallData)

	slog.Info("starting server", slog.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func decodeContractCallData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req communicator.DecodeContractCallDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, "Failed to decode request body", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respStruct, err := communicator.DecodeContractCallData(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(respStruct)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal response", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func getTransactionByHash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	respStruct, err := communicator.GetTransactionByHash(ctx, communicator.GetTransactionByHashRequest{
		Hash: chi.URLParam(r, "hash"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(respStruct)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal response", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	numberOfBlocks := r.URL.Query().Get("number_of_blocks")
	numberOfBlocksInt, err := strconv.Atoi(numberOfBlocks)
	if err != nil || numberOfBlocksInt <= 0 {
		numberOfBlocksInt = 3 // Default to 3 blocks if not provided or invalid
	}

	blockNumber := r.URL.Query().Get("block_number")
	blockNumberInt, err := strconv.Atoi(blockNumber)
	if err != nil || blockNumberInt < 0 {
		blockNumberInt = 0
	}

	respStruct, err := communicator.GetLatestNBlock(ctx, communicator.GetLatestNBlockRequest{
		NumberOfBlocks: int64(numberOfBlocksInt),
		BlockNumber:    int64(blockNumberInt),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(respStruct)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal response", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", slog.Any("err", err))
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
