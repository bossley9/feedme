PREFIX = /usr/local
BIN = $(PREFIX)/bin

EXE = feedme

build: test
	go build -o ./$(EXE) ./cmd/feedme

test:
	go test ./pkg/atom

clean:
	rm -f ./$(EXE)

server:
	go run cmd/feedme/main.go

install: build
	cp ./$(EXE) $(BIN)/$(EXE)
	chmod 555 $(BIN)/$(EXE)

uninstall:
	rm -f $(BIN)/$(EXE)

.PHONY: build test clean server install uninstall
