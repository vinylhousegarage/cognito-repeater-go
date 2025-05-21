.PHONY: format lint prepare

format:
	goimports -w .

lint:
	golangci-lint run

prepare: format lint
