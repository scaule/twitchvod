.PHONY: build build-cli build-server run-cli run-server test clean

# Build both CLI and server
build: build-cli build-server

# Build CLI version
build-cli:
	@echo "Building CLI version..."
	go build -o bin/twitchvod cmd/cli/main.go

# Build HTTP server version
build-server:
	@echo "Building HTTP server version..."
	go build -o bin/twitchvod-server cmd/server/main.go

# Run CLI version (requires TWITCH_CLIENT_ID env var)
run-cli: build-cli
	@echo "Running CLI version..."
	@echo "Usage: ./bin/twitchvod <twitch_url>"
	@echo "Make sure TWITCH_CLIENT_ID environment variable is set"

# Run HTTP server (requires TWITCH_CLIENT_ID and SECRET env vars)
run-server: build-server
	@echo "Running HTTP server..."
	@echo "Make sure TWITCH_CLIENT_ID and SECRET environment variables are set"
	@echo "Server will start on port 8080 (or PORT env var)"
	@echo "Usage: GET /extract?url=<twitch_url>&secret=<your_secret>"
	./bin/twitchvod-server

# Run tests
test:
	go test ./internal/twitch -v

# Clean build artifacts
clean:
	rm -rf bin/

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build both CLI and server"
	@echo "  build-cli   - Build CLI version only"
	@echo "  build-server- Build HTTP server version only"
	@echo "  run-cli     - Build and show CLI usage"
	@echo "  run-server  - Build and run HTTP server"
	@echo "  test        - Run tests"
	@echo "  clean       - Remove build artifacts"
	@echo "  help        - Show this help" 