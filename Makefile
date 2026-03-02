.PHONY: test fmt vet run-api

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	@echo "lint: add golangci-lint later (E0/E1)"

run-api:
	go run ./cmd/api