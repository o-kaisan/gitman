.PHONY: run
run:
	go run gitman.go $(ARGS)

.PHONY: test
test:
	go test -p 2 ./...

.PHONY: clean
clean:
	rm -f gm

.PHONY: build
build:
	make clean
	go build gitman.go
	mv ./gitman ./gm

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: install
install:
	bash install.sh
