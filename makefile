all: build
build:
	go build -o bin/emm-local-proxy -v
	cp publickey bin/publickey