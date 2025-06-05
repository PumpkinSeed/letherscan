# Stage 1: Build the frontend
FROM node:23 AS frontend-builder

WORKDIR /app
COPY frontend ./frontend
RUN npm --prefix ./frontend install
RUN npm --prefix ./frontend run build

# Stage 2: Build the Go backend
FROM golang:1.23 AS backend-builder

WORKDIR /app
COPY . .

# Copy frontend build artifacts from frontend-builder
COPY --from=frontend-builder /app/frontend/build ./bin/build

# Build the backend
RUN go build -o ./bin/letherscan ./bin/main.go

# Final stage: create a minimal image to run the backend
FROM debian:bookworm-slim

COPY --from=backend-builder /app/bin/letherscan /app/bin/letherscan

ENTRYPOINT ["/app/bin/letherscan"]