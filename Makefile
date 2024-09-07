FINDFILES=find . \( -path ./.git -o -path ./out -o -path ./.github -o -path ./vendor -o -path ./frontend/node_modules \) -prune -o -type f
XARGS=xargs -0 -r
RELEASE_LDFLAGS='-extldflags -static -s -w'
BINARIES=./cmd/bot
HUB ?= adhp
IMG ?= bot
TAG ?= dev

lint:
	@${FINDFILES} -name '*.go' \( ! \( -name '*.gen.go' -o -name '*.pb.go' \) \) -print0 | ${XARGS} scripts/lint_go.sh

.PHONY: default
default: init build

.PHONY: init
init:
	@mkdir -p out

.PHONY: build
build:
	@LDFLAGS=${RELEASE_LDFLAGS} scripts/gobuild.sh out/ ${BINARIES}

.PHONY: docker
docker:
	@./scripts/docker.sh

.PHONY: mod-vendor
mod-vendor:
	@go mod vendor

.PHONY: dev
dev:
	@go run ./cmd/bot/main.go --log-level debug start

.PHONY: clean
clean:
	@rm -rf out

.PHONY: dist-clean
dist-clean: clean
	@rm -rf vendor
