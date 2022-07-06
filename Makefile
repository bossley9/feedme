VERSION = v0.0.3

build:
	go build cmd/main.go

server:
	go run cmd/main.go

test:
	go test ./pkg/atom

publish:
	GOPROXY=proxy.golang.org go list -m git.sr.ht/~bossley9/feed-generator@$(VERSION)
