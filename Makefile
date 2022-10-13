run: build
	./lg

build: lg.go pkg/**/*
	go build -ldflags="-s -w" -o lg .

test: pkg/**/*
	go test -count=1 -v ./test

clean:
	rm ./lg

lint: install-golangci-lint
	golangci-lint run 

install-golangci-lint:
	which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint