#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

# Grab the app name from the bundler.json file
APP_NAME := $(shell scripts/get_app_name.py)

default: build

build:
	cd src/; astilectron-bundler -v

run: build
	./bin/linux-amd64/'${APP_NAME}'

run_debug: build
	./bin/linux-amd64/'${APP_NAME}' -d

webdev:
	./scripts/run_browsersync.sh

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./...

clean:
	rm -Rf bin/
