BINARY_NAME=minishell
BINARY_DIR=bin
CMD_DIR=cmd

.PHONY: all build clean test run help

all: build ## Собрать проект

build: ## Собрать бинарный файл
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BINARY_DIR)/$(BINARY_NAME)"

clean: ## Очистить артефакты сборки
	@go clean
	@rm -rf $(BINARY_DIR)

test: ## Запустить тесты
	@go test -v ./internal/parser
	@go test -v ./internal/executor

run: build ## Собрать и запустить
	@./$(BINARY_DIR)/$(BINARY_NAME)

help: ## Показать справку
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
