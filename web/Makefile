BIN_DIR := $(GOPATH)/bin
MUSTACHE := $(BIN_DIR)/mustache
NODE_MODULES := node_modules/

$(MUSTACHE):
	go get github.com/cbroglie/mustache/...

$(NODE_MODULES):
	yarn

.PHONY: gen

gen: $(MUSTACHE)
	@./gen_static.sh

INDEX_HTML?=src/index.html

$(INDEX_HTML):
	$(MAKE) gen

DIST?=dist/

$(DIST): $(INDEX_HTML)
	yarn build

.PHONY: build

build: $(DIST)

.PHONY: watch

watch: $(DIST)
	@yarn watch

.PHONY: clean

clean:
	@rm -rf .cache dist $(INDEX_HTML)s