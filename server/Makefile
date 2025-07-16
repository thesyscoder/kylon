# @File: Makefile
# @Title: Project Automation Makefile for Kylon Backend
# @Description: Provides a set of commands for building, running, testing, and cleaning
# @Description: the kylon backend application, facilitating development workflow automation.
# @Author: thesyscoder (github.com/thesyscoder)

# Project Variables
# Name of the executable binary that will be built
BINARY_NAME=kylon-server
# Directory where the compiled binary will be placed
BUILD_DIR=bin
# Path to the main Go package (the executable entry point)
CMD_PATH=./cmd

# Default target: runs the application with auto-reload for development.
# This makes 'make' equivalent to 'make run'.
.PHONY: all
all: run

# run: Starts the application using CompileDaemon for automatic recompilation and restart
# on code changes. Ideal for development.
# It sources environment variables from kylon_config.env and pipes output through jq for colorful logs.
.PHONY: run
run: build # Ensure the binary is built before running CompileDaemon
	@echo "🚀 Starting Kylon Backend with auto-reload and colorful logs..."
	# Check if kylon_config.env exists and load environment variables from it.
	# This makes sure settings like MinIO credentials are available to the Go app.
	@if [ -f .env ]; then \
		echo "Loading environment variables from .env..."; \
		export $(grep -v '^#' .env | xargs); \
	fi
	# CompileDaemon watches for file changes, rebuilds, and restarts the command.
	# The command itself pipes its stderr (where logrus writes) to stdout, then to jq for colorization.
	CompileDaemon \
		-build="go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)" \
		-command="$(BUILD_DIR)/$(BINARY_NAME) 2>&1 | jq ." \
		-directory="." \
		-include=".*\.go" \
		-exclude=".git" \
		-polling # Use polling for environments where file system events might not be reliable (e.g., some shared volumes)

# build: Compiles the Go application into an executable binary.
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME) binary..."
	# Ensure the build directory exists
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

# start: Runs the pre-built application binary manually.
# It sources the kylon_config.env file and pipes output through jq for colorful logs.
.PHONY: start
start: build # Ensure the binary is built before running
	@echo "🟢 Running $(BINARY_NAME) manually..."
	# Load environment variables just like in the 'run' target
	@if [ -f .env ]; then \
		echo "Loading environment variables from .env..."; \
		export $(grep -v '^#' .env | xargs); \
	fi
	# Execute the binary and pipe its output through jq for colorization
	$(BUILD_DIR)/$(BINARY_NAME) 2>&1 | jq .

# test: Executes all Go tests in the project.
.PHONY: test
test:
	@echo "🧪 Running tests..."
	go test ./...

# clean: Removes build artifacts, specifically the build directory.
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# setup-minio: Runs a local MinIO instance using Docker.
# This is useful for local development and testing of the storage integration.
.PHONY: setup-minio
setup-minio:
	@echo "🐳 Starting local MinIO instance (http://localhost:9001 for console)..."
	docker run --rm --name kylon-minio -p 9000:9000 -p 9001:9001 \
		-e MINIO_ROOT_USER=minioadmin \
		-e MINIO_ROOT_PASSWORD=minioadmin \
		quay.io/minio/minio server /data --console-address ":9001"

# stop-minio: Stops and removes the local MinIO Docker container.
.PHONY: stop-minio
stop-minio:
	@echo "🛑 Stopping local MinIO instance..."
	# '|| true' prevents the command from failing if the container doesn't exist
	docker stop kylon-minio || true && docker rm kylon-minio || true

# lint: Runs golangci-lint for code quality checks.
# Installs golangci-lint if not found.
.PHONY: lint
lint:
	@echo "🔍 Running linters..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...

# docker-build: Builds the Docker image for the backend.
# Assumes Dockerfile is in the backend/ directory.
.PHONY: docker-build
docker-build:
	@echo "📦 Building Docker image for $(BINARY_NAME)..."
	docker build -t thesyscoder:kylon:latest -f Dockerfile .

# docker-run: Runs the Docker image locally.
# Maps port 8080 (default server port) from container to host.
.PHONY: docker-run
docker-run: docker-build
	@echo "🏃 Running Docker image for $(BINARY_NAME) locally..."
	docker run --rm -p 8080:8080 --name kylon-local thesyscoder:kylon:latest
