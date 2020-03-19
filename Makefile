SHELL := /bin/bash

VERSION := $(CI_COMMIT_TAG)
GITCOMMIT := $(git rev-list -1 HEAD)
BRANCH := $(CI_COMMIT_BRANCH)
BUILDDATE := `date +%Y-%m-%d`
BUILDUSER := `whoami`
PROJECT_ROOT := github.com/benzcash/cloudlog

LDFLAGSSTRING :=-X main.Version=$(VERSION)
LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.Branch=$(BRANCH)
LDFLAGSSTRING +=-X main.BuildDate=$(BUILDDATE)
LDFLAGSSTRING +=-X main.BuildUser=$(BUILDUSER)

LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

.PHONY: all build

all: build

# Build binary
build:
	go build $(LDFLAGS) 

test:
	go test -v ./...