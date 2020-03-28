TERRAFORM_PROVIDER := terraform-provider-s3

GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

.PHONY: default
default: build

.PHONY: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(TERRAFORM_PROVIDER)

.PHONY: fmt
fmt:
	gofmt -s -w "$(CURDIR)/."

.PHONY: clean
clean:
	go clean
	go mod tidy
