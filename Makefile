.PHONY: build test clean install help

# Build the binary
build:
	go build -o manage-agent-skills .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f manage-agent-skills

# Install to /usr/local/bin (requires sudo)
install: build
	sudo mv manage-agent-skills /usr/local/bin/

# Run linting
lint:
	go fmt ./...
	go vet ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build   - Build the binary"
	@echo "  test    - Run tests"
	@echo "  clean   - Remove build artifacts"
	@echo "  install - Install to /usr/local/bin (requires sudo)"
	@echo "  lint    - Run linting"
	@echo "  help    - Show this help message"
