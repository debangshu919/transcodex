BIN_DIR = bin
API_BIN = $(BIN_DIR)/api
API_SRC = cmd/api/main.go

.PHONY: build run clean

build:
	@go build -o $(API_BIN) $(API_SRC)

run: build
	@./$(API_BIN)

clean:
	@rm -rf $(BIN_DIR)/
