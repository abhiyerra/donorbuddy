build: test
	go build

deps:
	go get -u

test:
	go test -v

dev-run: build
	./donorbuddy config.dev.json
