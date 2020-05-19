BIN_DIR := $(GOPATH)/bin
MUSTACHE := $(BIN_DIR)/mustache
NODE_MODULES := node_modules/

$(MUSTACHE):
	go get github.com/cbroglie/mustache/...

$(NODE_MODULES):
	yarn

gen: $(MUSTACHE)
	@./gen_static.sh

build: $(NODE_MODULES) src/index.html
	yarn build
