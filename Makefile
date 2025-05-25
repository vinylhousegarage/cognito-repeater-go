check-format:
	@unformatted=$$(goimports -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted files:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

format:
	goimports -w .
