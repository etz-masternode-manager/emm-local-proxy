all: build
build:
	go build -o bin/emm-proxy -v
	cp publickey bin/publickey