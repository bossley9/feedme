EXE = feedme

build: test
	./build.sh ./$(EXE)

test:
	go test ./pkg/atom

server:
	go run cmd/*.go
