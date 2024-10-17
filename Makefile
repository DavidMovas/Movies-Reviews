before-push:
	go mod tidy && \
	gofumpt -l -w . && \
	go build ./... && \
	go test ./integration_tests/...