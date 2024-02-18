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

BUILDDEPS := docs

ifeq ($(DOCS),false)
	BUILDDEPS :=
endif

docs:
	@echo "Generating swagger documentation"
	@swag init --pd -g server/server.go -d src --md src/docs -o src/docs

build: $(BUILDDEPS)
	@echo "Building $(OS) $(ARCH) binary..."
	@GOOS=$(OS) GOARCH=$(ARCH) go build $(ARGS) -o "bin/$(BINARY_NAME)" src/main.go

build_all: $(BUILDDEPS)
	@echo "Building linux amd64 binary..."
	@GOOS=linux GOARCH=amd64 go build -o "bin/$(BINARY_NAME)-linux-amd64" src/main.go
	@echo "Building darwin amd64 binary..."
	@GOOS=darwin GOARCH=amd64 go build -o "bin/$(BINARY_NAME)-darwin-amd64" src/main.go
	@echo "Building darwin arm64 binary..."
	@GOOS=darwin GOARCH=arm64 go build -o "bin/$(BINARY_NAME)-darwin-arm64" src/main.go
	@echo "Building windows amd64 binary..."
	@GOOS=windows GOARCH=amd64 go build -o "bin/$(BINARY_NAME)-windows-amd64.exe" src/main.go

test_clean:
	@echo "Cleaning test databases..."
	@rm -f app.db
	@rm -f events.db
	@rm -f events.db.wal
	@echo "Cleaning test results..."
	@rm -f coverage.unit.out
	@rm -f coverage.e2e.out

test_unit: test_clean
	@echo "Running unit tests..."
	@@GOOS=$(OS) GOARCH=$(ARCH) ENV=test go test -count 1 -timeout 10s -race -coverprofile=coverage.unit.out -covermode=atomic -v -coverpkg=./src/...  ./src/...

test_e2e: test_clean
	@echo "Running end-to-end tests..."
	@@GOOS=$(OS) GOARCH=$(ARCH) ENV=test go test -count 1 -timeout 10s -race -coverprofile=coverage.e2e.out -covermode=atomic -v -coverpkg=./src/...  ./tests/...

test: test_unit test_e2e

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf bin/*

run: build
	@echo "Running..."
	@./bin/$(BINARY_NAME)