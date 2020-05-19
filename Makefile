PROJECT?=github.com/anatollupacescu/retail-sample

GOOS?=linux
GOARCH?=amd64

RELEASE?=0.0.0
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

format:
	@goimports -w -local "github.com/anatollupacescu/retail-sample" cmd/ internal/

test:
	@go test ./...

build/docker:
	@docker build .

build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-ldflags "-s -w -X ${PROJECT}/internal/version.Version=${RELEASE} \
		-X ${PROJECT}/internal/version.Commit=${COMMIT} \
		-X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
		-o bin/retail ${PROJECT}/cmd/retail-sample

run:
	@go run $(shell pwd)/cmd/retail-sample

BINARY?=$(shell pwd)/bin/retail

$(BINARY):
	$(MAKE) build

.PHONY: run/mem

run/mem: $(BINARY)
	$(BINARY) --in-memory

.PHONY: run/docker

run/docker: $(BINARY)
	docker-compose up --build

.PHONY: clean

clean:
	@rm $(BINARY)
	@$(MAKE) clean -C web/

.PHONY: build/web

build/web:
	$(MAKE) build -C web/
 
.PHONY: test format build build/docker run
