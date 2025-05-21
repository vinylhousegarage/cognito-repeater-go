.PHONY: format lint

format:
	goimports -w .

lint:
	golangci-lint run
