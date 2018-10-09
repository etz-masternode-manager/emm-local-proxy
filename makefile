all: clean build
clean:
	rm -rf bin
build:
	go build -o bin/emm-proxy -v
	cp publickey bin/publickey