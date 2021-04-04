VERSION=0.1.0
REVISION=`git rev-parse HEAD`

export GO111MODULE=on

all: build

lint:
	@golint -min_confidence 0.3 -set_exit_status ./mux

build:
	@go build --buildmode=pie -o bin/monitor -ldflags "-s -w -X main.revision=${REVISION} -X main.version=${VERSION}"

clean:
	@rm -rf ./bin
