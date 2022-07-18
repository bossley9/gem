PREFIX = /usr/local
BIN = $(PREFIX)/bin
EXE = gem

make: test
	go build -o ./$(EXE) cmd/main.go

test:
	go test

clean:
	rm ./$(EXE)

install:
	mkdir -p $(BIN)
	cp -f ./$(EXE) $(BIN)
	chmod 555 $(BIN)/$(EXE)

uninstall:
	rm -f $(BIN)/$(EXE)

.PHONY: make test clean install uninstall
