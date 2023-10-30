SHELL:=/bin/bash
folder=$(shell pwd)
builtOn=$(shell hostname)
builtAt=$(shell date +'%F %T %z')
goVersion=$(shell go version | sed 's/go version //')
gitAuthor=$(shell git show -s --format='format:%aN <%ae>' HEAD)
gitCommit=$(shell git rev-parse HEAD)--$(shell git -C "$(folder)" rev-parse --abbrev-ref HEAD)
gitTag:=$(shell git -C "$(folder)" describe --tags HEAD|| git -C "$(folder)" rev-parse --abbrev-ref HEAD | grep -v HEAD || git -C "$(folder)" rev-parse HEAD )
GOLDFLAGS = -s -w -X 'go-lib/utils/build/inject.builtOn=$(builtOn)'
GOLDFLAGS += -X 'go-lib/utils/build/inject.builtAt=$(builtAt)'
GOLDFLAGS += -X 'go-lib/utils/build/inject.builtUser=$(shell whoami)'
GOLDFLAGS += -X 'go-lib/utils/build/inject.goVersion=$(goVersion)'
GOLDFLAGS += -X 'go-lib/utils/build/inject.gitAuthor=$(gitAuthor)'
GOLDFLAGS += -X 'go-lib/utils/build/inject.gitCommit=$(gitCommit)'
GOLDFLAGS += -X 'go-lib/utils/build/inject.gitTag=$(gitTag)'

.PHONY: build
build:
	go build -o inject.exe -ldflags "$(GOLDFLAGS)" main.go && ./inject.exe -version