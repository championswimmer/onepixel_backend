BINARY_NAME := "onepixel"

ARCH := $(or $(GOARCH),$(shell uname -m))
OS := $(or $(GOOS),$(shell uname))

ifneq (,$(filter $(OS),Darwin darwin MacOS macos))
	OS := darwin
else ifneq (,$(filter $(OS),Linux linux))
	OS := linux
else
	OS := windows
endif

ifeq ($(ARCH),x86_64)
	ARCH := amd64
else ifeq ($(ARCH),i386)
	ARCH := 386
endif

docs:
	@echo "Generating swagger documentation"
	@swag init -g src/main.go -o src/docs

build: docs
	@echo "Building $(OS) $(ARCH) binary..."
	@GOOS=$(OS) GOARCH=$(ARCH) go build $(ARGS) -o "bin/$(BINARY_NAME)" src/main.go

build_all: docs
	@echo "Building linux amd64 binary..."
	@GOOS=linux GOARCH=amd64 go build -o "bin/$(BINARY_NAME)-linux-amd64" src/main.go
	@echo "Building darwin amd64 binary..."
	@GOOS=darwin GOARCH=amd64 go build -o "bin/$(BINARY_NAME)-darwin-amd64" src/main.go
	@echo "Building darwin arm64 binary..."
	@GOOS=darwin GOARCH=arm64 go build -o "bin/$(BINARY_NAME)-darwin-arm64" src/main.go
	@echo "Building windows amd64 binary..."
	@GOOS=windows GOARCH=amd64 go build -o "bin/$(BINARY_NAME)-windows-amd64.exe" src/main.go

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf bin/*

run: build
	@echo "Running..."
	@./bin/$(BINARY_NAME)-$(OS)-$(ARCH)