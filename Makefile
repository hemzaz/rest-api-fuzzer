PROJECT_NAME := rest-api-fuzzer
BINARY_NAME := fuzzer

GO_PACKAGES := github.com/spf13/cobra github.com/go-resty/resty/v2 go.etcd.io/bbolt

.PHONY: all
all: init build

.PHONY: init
init:
	@echo "Initializing project..."
	@go mod init $(PROJECT_NAME) || true
	@go get -u $(GO_PACKAGES)

.PHONY: build
build:
	@echo "Building the project..."
	@go build -o $(BINARY_NAME)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f responses.db

.PHONY: probe
probe:
	@echo "Running in probe mode..."
	./$(BINARY_NAME) probe --url http://localhost:8080

.PHONY: breach
breach:
	@echo "Running in breach mode..."
	./$(BINARY_NAME) breach --url http://localhost:8080 --threads 8

.PHONY: get-packages
get-packages:
	@echo "Getting required Go packages..."
	@go get -u $(GO_PACKAGES)

.PHONY: fmt
fmt:
	@echo "Formatting the code..."
	@go fmt ./...

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: full
full: clean init build
