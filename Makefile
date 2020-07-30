goversion=1.14.6
golangci_lint_version=1.29
vols=-v `pwd`:/app -w /app
run_go=docker run --rm $(vols) golang:$(goversion)
run_lint=docker run --rm $(vols) golangci/golangci-lint:v$(golangci_lint_version)
cov=coverage.out
covhtml=coverage.html

.PHONY: all
all: test lint

.PHONY: test
test:
	@$(run_go) go test -coverprofile=$(cov) -race ./...

.PHONY: coverage
coverage: test
	@$(run_go) go tool cover -html=$(cov) -o=$(covhtml)
	@open $(covhtml) || xdg-open $(covhtml)

.PHONY: lint
lint:
	@$(run_lint) run ./...
