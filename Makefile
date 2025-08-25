.PHONY: run
run:
	go run gitman.go $(ARGS)

.PHONY: test
test:
	go test -p 2 ./...

.PHONY: clean
clean:
	rm -f gitman

.PHONY: build
build:
	make clean
	go build gitman.go

.PHONY: lint
lint:
	staticcheck ./...
