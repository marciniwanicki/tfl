.PHONY: build test lint format clean

build:
	go build -o tfl .

test:
	go test -v ./...

lint:
	golangci-lint run

format:
	go fmt ./...

clean:
	rm -f tfl
