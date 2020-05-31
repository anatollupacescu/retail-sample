PROJECT:=$(shell go list -m)

GOOS?=linux
GOARCH?=amd64

RELEASE?=0.0.0
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

format:
	@goimports -w -local $(PROJECT) cmd/ internal/

test:
	@go test ./...

build/docker:
	@docker build .

build:
	echo $(PROJECT)
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
	$(MAKE) build/web
	docker-compose up --build

.PHONY: clean/api clean/web

clean/api:
	@rm $(BINARY) 2> /dev/null || true
	# @$(MAKE) clean -C web/

clean/web:
	@$(MAKE) clean -C web/

.PHONY: build/web

build/web:
	$(MAKE) build -C web/
 
.PHONY: test format build build/docker run

BIN_DIR := $(shell go env GOPATH)/bin
GRAPH_TOOL := $(BIN_DIR)/godepgraph

$(GRAPH_TOOL):
	go get github.com/kisielk/godepgraph

.PHONY: graph

graph: $(GRAPH_TOOL)
	@$(GRAPH_TOOL) -s -novendor \
		-o ./cmd,$(PROJECT) \
		./cmd/retail-sample/ | dot -Tpng -o graph.png