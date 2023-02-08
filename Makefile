default: tool fmt tidy lint

tool:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/daixiang0/gci@latest

.PHONY: migrate-test
migrate-test:
	@goose --dir="deployments/migrations" postgres "postgres://pg-user:pg-pass@127.0.0.1:5436/pg-db?sslmode=disable" down
	@goose --dir="deployments/migrations" postgres "postgres://pg-user:pg-pass@127.0.0.1:5436/pg-db?sslmode=disable" up

.PHONY: migrate-run
migrate-run:
	@goose --dir="deployments/migrations" postgres "postgres://pg-user:pg-pass@127.0.0.1:5436/pg-db?sslmode=disable" up

.PHONY: lint
lint: $(GOLANGCI_BIN)
	@golangci-lint run

.PHONY: test
test:
	@go test ./... -cover -race

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: fmt
fmt: $(GOFUMPT_BIN) $(GCI_BIN)
	@gofumpt -l -w .
	@gci write -s standard -s default --skip-generated .
