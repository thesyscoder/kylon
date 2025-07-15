# @File: Makefile
# @Title: Project Automation Makefile
# @Description: Provides a set of commands for building, running, testing, and cleaning
# @Description: the kylon application, facilitating development workflow automation.
# @Author: thesyscoder (github.com/thesyscoder)

# Project Variables
BINARY_NAME=kylon
BUILD_DIR=bin
CMD_DIR=./cmd

# Default target: runs the application with auto-reload for development.
.PHONY: all
all: run

# run: Starts the application using CompileDaemon for automatic recompilation and restart
# on code changes. Ideal for development.
.PHONY: run
run:
	@echo "🚀 Starting CompileDaemon with auto-reload..."
	CompileDaemon \
		-build="go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)" \
		-command="$(BUILD_DIR)/$(BINARY_NAME)" \
		-directory="." \
		-include=".*\.go" \
		-exclude=".git" \
		-polling

# build: Compiles the Go application into an executable binary.
.PHONY: build
build:
	@echo "🔨 Building binary..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# start: Runs the pre-built application binary manually.
.PHONY: start
start:
	@echo "🟢 Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

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
