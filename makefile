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
	@trap 'kill %1; kill %2' SIGINT; \
	cd frontend && (npm run dev 2>&1 | sed 's/^/\x1b[36m[Frontend]\x1b[0m /') & \
	(go run main.go 2>&1 | sed 's/^/\x1b[32m[Backend]\x1b[0m  /') & \
	wait