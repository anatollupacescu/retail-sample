BIN_DIR := $(GOPATH)/bin
TERN := $(BIN_DIR)/tern

$(TERN):
	@go get -u github.com/jackc/tern

migrate: $(TERN)
	@tern migrate

.PHONY: migrate
