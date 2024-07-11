# Build Stage
FROM golang:1.22.2-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo-service cmd/main.go

# Run Stage
FROM alpine as runner

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/todo-service ./app

ENTRYPOINT ["./app"]
