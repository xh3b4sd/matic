PROJECT=simple
ORGANIZATION=zyndiecate
SOURCE=$(shell find . -name '*.go')

GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)

.PHONY=all clean test

all: $(GOPATH) $(PROJECT)

clean:
	rm -rf $(GOPATH) $(PROJECT)

# deps
$(GOPATH):
	mkdir -p $(PROJECT_PATH)
	cd $(PROJECT_PATH) && ln -s ../../../.. $(PROJECT)

	#
	# Fetch private packages first (so `go get` skips them later)

	#
	# Fetch public dependencies via `go get`
	GOPATH=$(GOPATH) go get -d github.com/$(ORGANIZATION)/$(PROJECT)

	#
	# Build test packages (we only want those two, so we use `-d` in go get)
	GOPATH=$(GOPATH) go get github.com/onsi/gomega
	GOPATH=$(GOPATH) go get github.com/onsi/ginkgo

# build
$(PROJECT): $(SOURCE)
	GOPATH=$(GOPATH) go build -o $(PROJECT)
