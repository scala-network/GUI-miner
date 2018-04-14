#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := stellite-api

default: build

build:
	go build -o ./bin/${APP_NAME} ./src/main.go

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./src/...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

run: build
	LOG_FORMAT=Text \
	LOG_LEVEL=Debug \
	DATABASE_ENDPOINT="127.0.0.1:3306" \
	DATABASE_NAME="stellite_chain" \
	DATABASE_USERNAME="root" \
	DATABASE_PASSWORD="root" \
	XTL_DAEMON_ENDPOINT="127.0.0.1:20189" \
	EXCHANGE_TRADEOGRE_ENDPOINT="https://tradeogre.com/api/v1/ticker/%s-%s" \
	EXCHANGE_CREX24_ENDPOINT="https://api.crex24.com/CryptoExchangeService/BotPublic/ReturnTicker?request=[NamePairs=%s_%s]" \
	POOL_SLEEP_SECONDS=5 \
	EXCHANGE_SLEEP_SECONDS=5 \
	BLOCKCHAIN_SLEEP_SECONDS=5 \
	./bin/${APP_NAME}

run_race:
	LOG_OUTPUT=Text \
	LOG_LEVEL=Debug \
	DATABASE_ENDPOINT=127.0.0.1:3306 \
	DATABASE_NAME=stellite_chain \
	DATABASE_USERNAME=root \
	DATABASE_PASSWORD=root \
	XTL_DAEMON_ENDPOINT=127.0.0.1:20189 \
	go run -race ./src/main.go

clean:
	rm ./bin/*
