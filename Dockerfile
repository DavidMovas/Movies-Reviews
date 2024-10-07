FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.sum go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:latest AS final

COPY --from=builder /app/server /app/server

ENTRYPOINT ["/app/server"]