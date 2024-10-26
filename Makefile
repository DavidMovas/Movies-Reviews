run:
	go run .

compose:
	docker-compose up --build -d

go-lint:
	golangci-lint run ./...

go-fmt:
	 gofumpt -l -w .

tidy:
	go mod tidy

test:
	go test ./integration_tests/...

before-push:
	go mod tidy && \
	gofumpt -l -w . && \
	go build ./... && \
	golangci-lint run ./... && \
	go test ./integration_tests/...