Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

run-only:
	./build/debug/wizard

run: build run-only

run-dashboard:
	cd ./dashboard && npm run dev

build-dashboard:
	cd ./dashboard && npm run build

static-gen: build-dashboard
	esc -pkg api -o api/assets.go -prefix=dashboard/dist dashboard/dist
	go generate ./cmd/wizard

build:
	go build -race -ldflags $(LDFLAGS) -o build/debug/wizard cmd/web/main.go

build-lorca:
	go build -race -ldflags $(LDFLAGS) -o build/debug/wizard cmd/lorca/main.go

build-systray:
	go build -race -ldflags $(LDFLAGS) -o build/debug/wizard cmd/systray/main.go

clean:
	rm -fr build/debug/wizard build/release/wizard*

.PHONY: run build build-release clean
