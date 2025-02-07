.PHONY: build

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

build:
	@cd frontend && npm run build
	@echo "Building API documentation"
	@npx -y @redocly/cli build-docs erugo.openapi.json --output=swagger.html
	@echo "Building for GOOS=$(GOOS) GOARCH=$(GOARCH)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o erugo-$(GOOS)-$(GOARCH)
	@echo "✓ Binary built successfully: erugo-$(GOOS)-$(GOARCH)"

run:
	@cd frontend && npm run dev &
	@go run main.go