.PHONY: build run test clean

build:
	go build -o towncrier cmd/towncrier/main.go

run: build
	./towncrier

test:
	go test ./...

clean:
	rm -rf towncrier