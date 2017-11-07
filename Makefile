BIN      := docker_stats 
OSARCH   := "linux/amd64" 
VERSION  := $(shell git describe --tags)

all: build

test: build
	./test.sh

deps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get github.com/mitchellh/gox

lint: deps
	go vet ./...
	golint -set_exit_status ./...

package:
	rm -fR ./pkg && mkdir ./pkg ;\
		gox \
		-osarch $(OSARCH) \
		-output "./pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" \
		-ldflags "-X main.version=$(VERSION)" \
		.;\
	    for d in $$(ls ./pkg);do zip ./pkg/$${d}.zip ./pkg/$${d}/*;done

build:
	go build -o $(BIN) -ldflags "-X main.version=$(VERSION)" .

clean:
	go clean
