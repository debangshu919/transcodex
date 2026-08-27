BIN_DIR = bin
API_BIN = $(BIN_DIR)/api
API_SRC = cmd/api/main.go

.PHONY: build run clean migrate-up migrate-down

build:
	@go build -o $(API_BIN) $(API_SRC)

run: build
	@./$(API_BIN)

clean:
	@rm -rf $(BIN_DIR)/

STEPS ?=

migrate-up:
ifeq ($(STEPS),)
	@go run cmd/migrate/main.go up
else
	@go run cmd/migrate/main.go up $(STEPS)
endif

migrate-down:
ifeq ($(STEPS),)
	@go run cmd/migrate/main.go down
else
	@go run cmd/migrate/main.go down $(STEPS)
endif
