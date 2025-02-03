.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	golangci-lint run ./... --timeout 5m

.PHONY: test
test:
	go test ./... -gcflags 'all=-l' -failfast -timeout 20s -coverprofile .test-coverage.txt

.PHONY: coverage-report
coverage-report: test
	go tool cover -func=.test-coverage.txt

.DEFAULT_GOAL := test
