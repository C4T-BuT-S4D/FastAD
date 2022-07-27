.PHONY: lint-proto
lint-proto:
	cd proto && buf lint && buf build

.PHONY: lint-go
lint-go:
	golangci-lint run -v -c .golangci.yml ./...

.PHONY: proto
proto: lint-proto
	cd proto && buf generate

.PHONY: tidy
tidy:
	go mod tidy -compat="1.18"

.PHONY: lint
lint: lint-proto lint-go
