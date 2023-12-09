format:
	@gofumpt -l -w -extra .
lint:
	@golangci-lint run
