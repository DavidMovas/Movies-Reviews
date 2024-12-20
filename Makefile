run:
	go run .

compose:
	docker-compose up --build -d

compose-down:
	docker-compose down

go-lint:
	golangci-lint run ./...

go-fmt:
	 gofumpt -l -w .

gen:
	swag init

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