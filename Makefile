PROJECT:=$(shell go list -m)

GOOS?=linux
GOARCH?=amd64

RELEASE?=0.0.0
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: format test

format:
	@goimports -w -local $(PROJECT) cmd/ internal/

test:
	@go test ./...

# build

.PHONY: build build/web build/api build/docker

build/web:
	$(MAKE) build -C web/
 
build/docker:
	@docker build .

build/api:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-ldflags "-s -w -X ${PROJECT}/internal/version.Version=${RELEASE} \
		-X ${PROJECT}/internal/version.Commit=${COMMIT} \
		-X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
		-o bin/retail ${PROJECT}/cmd/retail-sample

build: build/web build/api

# run

.PHONY: run/mem run/docker

run:
	@go run $(shell pwd)/cmd/retail-sample

BINARY?=$(shell pwd)/bin/retail

$(BINARY):
	$(MAKE) build

run/mem: $(BINARY)
	$(BINARY) --in-memory

run/docker: $(BINARY)
	$(MAKE) build/web
	docker-compose up --build

# clean

.PHONY: clean clean/api clean/web

clean/api:
	@rm $(BINARY) 2> /dev/null || true

clean/web:
	@$(MAKE) clean -C web/

clean: clean/api clean/web

BIN_DIR := $(shell go env GOPATH)/bin
GRAPH_TOOL := $(BIN_DIR)/godepgraph

$(GRAPH_TOOL):
	go get github.com/kisielk/godepgraph

.PHONY: graph

graph: $(GRAPH_TOOL)
	@$(GRAPH_TOOL) -s -novendor \
		-o ./cmd,$(PROJECT) \
		./cmd/retail-sample/ | dot -Tpng -o graph.png

.PHONY: start/arbor

start/arbor:
	@go run $(shell pwd)/cmd/arbor

.PHONY: test/acceptance

test/acceptance:
	@go test $(shell pwd)/cmd/retail-sample-test/... \
	-v -tags=acceptance -args \
	--apiURL=http://localhost:8080/inventory \
	--arborURL=http://localhost:3000/data.json 
