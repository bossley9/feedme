PREFIX = /usr/local
BIN = $(PREFIX)/bin
RCBIN = /etc/rc.d

EXE = feedme
RC = feedme.rc

build: test
	./build.sh ./$(EXE)

test:
	go test ./pkg/atom

clean:
	rm -f ./$(EXE)

server:
	go run cmd/*.go

install: build
	cp ./$(EXE) $(BIN)/$(EXE)
	chmod 555 $(BIN)/$(EXE)
	cp ./$(RC) $(RCBIN)/$(EXE)
	chmod 755 $(RCBIN)/$(EXE)

uninstall:
	rm -f $(BIN)/$(EXE)
	rm -f $(RCBIN)/$(EXE)

.PHONY: build clean server install uninstall
