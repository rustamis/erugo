.PHONY: build run

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GIT_TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_VERSION := $(shell \
	TAG_COMMIT=$$(git rev-list -n1 $(GIT_TAG) 2>/dev/null); \
	HEAD_COMMIT=$$(git rev-parse HEAD 2>/dev/null); \
	if [ "$$TAG_COMMIT" = "$$HEAD_COMMIT" ]; then \
		echo "$(GIT_TAG)"; \
	else \
		echo "$(GIT_TAG)*"; \
	fi)

build:
	@cd frontend && npm run build
	@echo "Adding version tag $(GIT_VERSION) to index.html"
	@if [ -f frontend/dist/index.html ]; then \
		perl -pi -e 's/<body([^>]*)>/<body$$1 data-version="$(GIT_VERSION)">/' frontend/dist/index.html || \
		(echo "Failed to update index.html" && exit 1); \
	else \
		echo "Warning: frontend/dist/index.html not found"; \
	fi
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