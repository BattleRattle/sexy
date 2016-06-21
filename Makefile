BINARY = sexy
VERSION ?= 0.0.1
BUILD_TIME = $(shell date +%FT%T%z)
REVISION = $(shell git describe --dirty --always --tags)
LDFLAGS = "
LDFLAGS += -X github.com/BattleRattle/sexy/version.Version=$(VERSION)
LDFLAGS += -X github.com/BattleRattle/sexy/version.BuildTime=$(BUILD_TIME)
LDFLAGS += -X github.com/BattleRattle/sexy/version.CommitHash=$(REVISION)
LDFLAGS += "

.DEFAULT_GOAL: build
.PHONY: build

build:
	go build -ldflags $(LDFLAGS) -o $(BINARY) main.go
