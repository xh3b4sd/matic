PROJECT=matic
ORGANIZATION=zyndiecate

SOURCE := $(shell find . -name '*.go')
VERSION := $(shell cat VERSION)
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)

.PHONY=all clean test deps bin

all: deps bin

clean:
	rm -rf $(GOPATH) $(PROJECT) simple

test:
	GOPATH=$(GOPATH) go test ./...

# deps
deps: .gobuild
.gobuild:
	mkdir -p $(PROJECT_PATH)
	cd $(PROJECT_PATH) && ln -s ../../../.. $(PROJECT)

	#
	# Fetch private packages

	#
	# Fetch public packages
	GOPATH=$(GOPATH) go get -d github.com/$(ORGANIZATION)/$(PROJECT)
	GOPATH=$(GOPATH) go get -d github.com/$(ORGANIZATION)/$(PROJECT)/fixture/simple

	#
	# Fetch test packages
	GOPATH=$(GOPATH) go get -d github.com/onsi/gomega
	GOPATH=$(GOPATH) go get -d github.com/onsi/ginkgo

# build
bin: $(SOURCE)
	GOPATH=$(GOPATH) go build -ldflags "-X main.clientMaticVersion $(VERSION)" -o $(PROJECT)
	GOPATH=$(GOPATH) go build ./fixture/simple/...
