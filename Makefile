build:
	go build cmd/main.go

server:
	go run cmd/main.go

test:
	go test ./pkg/atom
