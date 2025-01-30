.PHONY: build

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

build:
	cd frontend && npm run build
	echo "Building for GOOS=$(GOOS) GOARCH=$(GOARCH)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/erugo-$(GOOS)-$(GOARCH)

run:
	cd frontend && npm run dev &
	go run main.go