goversion=1.14.6
golangci_lint_version=1.29
vols=-v `pwd`:/app -w /app
run_go=docker run --rm $(vols) golang:$(goversion)
run_lint=docker run --rm $(vols) golangci/golangci-lint:v$(golangci_lint_version)
gotests=go test -coverprofile=$(cov) -race ./...
cov=coverage.out
covhtml=coverage.html

export goversion

.PHONY: all
all: test lint

.PHONY: test
test:
	$(run_go) $(gotests)

# WHY: Had a lot of issues with stale containers and images
# not getting rebuilt on docker-compose. Also when test fails
# cleanup is not done (restoring db to initial state), so we
# always guarantee cleanup before running tests.
.PHONY: test-integration
test-integration: cleanup
	docker-compose build && \
	docker-compose run dev $(gotests) -tags integration && \
	docker-compose down -v

.PHONY: run
run: cleanup
	docker-compose build && \
	docker-compose run user-manager && \
	docker-compose down -v

.PHONY: coverage
coverage: test-integration
	@$(run_go) go tool cover -html=$(cov) -o=$(covhtml)
	@open $(covhtml) || xdg-open $(covhtml)

.PHONY: lint
lint:
	@$(run_lint) golangci-lint run ./...

.PHONY: shell
shell:
	docker-compose build && \
	docker-compose run dev && \
	docker-compose down

.PHONY: cleanup
cleanup:
	docker-compose down -v
