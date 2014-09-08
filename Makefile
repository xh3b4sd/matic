PROJECT=matic
ORGANIZATION=zyndiecate

SOURCE := $(shell find . -name '*.go')
VERSION := $(shell cat VERSION)
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)
FIXTURE_PATH := $(PROJECT_PATH)/$(PROJECT)/src/fixture

.PHONY=all clean test

all: $(GOPATH) $(PROJECT)

clean:
	rm -rf $(GOPATH) $(PROJECT)

test:
	GOPATH=$(GOPATH) go test ./...

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
	GOPATH=$(GOPATH) go build -ldflags "-X main.clientMaticVersion $(VERSION)" -o $(PROJECT)

build-fixture:
	cd $(FIXTURE_PATH)/simple && make

clean-fixture:
	cd $(FIXTURE_PATH)/simple/ && rm -rf simple .gobuild
