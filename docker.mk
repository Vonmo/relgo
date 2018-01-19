#!/usr/bin/env make -rRf

ifndef VERBOSE
MAKEFLAGS += --no-print-directory
endif

RELNAME:=relgo
RELEASE?=$(shell date +"%Y.%m")
COMMIT?=$(shell git rev-parse --short HEAD)
RELEASE_NAME:="${RELEASE}-${COMMIT}"
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILDER_NAME?=$(shell whoami)
BUILD_MACHINE?=$(shell hostname)


# use to override vars for your platform
ifeq (env.mk,$(wildcard env.mk))
include env.mk
endif

fmt:
	@gofmt -w ./

test: fmt migrate stop_prev
	@go test ./tests -v

prod: build optimize

build:
	@echo "Building..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 packr build -p 2 -ldflags "-s -w -X github.com/Vonmo/relgo/config.Version=${RELEASE} -X github.com/Vonmo/relgo/config.Commit=${COMMIT} -X github.com/Vonmo/relgo/config.BuildTime=${BUILD_TIME} -X github.com/Vonmo/relgo/config.BuilderName=${BUILDER_NAME} -X github.com/Vonmo/relgo/config.BuildMachine=${BUILD_MACHINE}" -o ./bin/${RELNAME}
	@chmod a+rwx ./bin/${RELNAME}

optimize:
	@echo "Optimization..."
	@upx -q ./bin/*
	@chmod a+rwx ./bin/*

run: build migrate stop_prev
	@./bin/${RELNAME} -config ./config/config.yml

dep:
	@dep ensure -v

new_migration:
	@sql-migrate new --config ./config/migrations.yml $n && chmod a+rw ./migrations/*
migrate:
	@sql-migrate up --config ./config/migrations.yml

stop_prev:
	@pkill ${RELNAME} || true