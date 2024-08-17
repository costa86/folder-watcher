.DEFAULT_GOAL := run
BINARY=folder-watcher
BUILD_CMD := go build --trimpath -ldflags="-s -w"
BUILD_CMD_LINUX := go build -ldflags="-s -w"
ARCH := amd64

run:
	@go run .

build:
	@go mod tidy

	@echo "Building for Linux"
	@GOOS=linux GOARCH=$(ARCH) $(BUILD_CMD_LINUX) -o $(BINARY)-linux

	@echo "Building for Windows" 
	@GOOS=windows GOARCH=$(ARCH) $(BUILD_CMD) -o $(BINARY)-windows.exe

	@echo "Building for MacOS"
	@GOOS=darwin GOARCH=$(ARCH) $(BUILD_CMD) -o $(BINARY)-darwin
