# Copyright 2019 VMware, Inc.
# SPDX-License-Identifier: BSD-2

# tag is the git tag or diff from the tag
TAG := $(shell git --no-pager describe --tags --always)

# or if empty, tag is the git revision
ifeq ($(TAG),)
TAG := $(shell git rev-parse --verify HEAD)
endif

# or if empty, tag is latest
ifeq ($(TAG),)
TAG := latest
endif

BASEREPO   ?= index.docker.io/tompscanlan
UIIMG      := incident-reporting-ui:${TAG}
TRUFFLEIMG := incident-reporting-truffle:${TAG}
LATESTTAG  := latest

GANACHE_CLI := ganache-test-incident-reporting

all: containers
containers: ui truffle
push: push-ui push-truffle
.PHONY: ui truffle push-ui push-truffle containers

ui: ui-server/*.go ui-server/Dockerfile truffle/contracts/*.sol
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose build ui

truffle: truffle/package*.json truffle/Dockerfile truffle/contracts/*.sol truffle/migrations/*.js truffle/test/*.js truffle/*.js
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose build truffle

push-ui: ui
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose push ui

push-truffle: truffle
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose push truffle

run-truffle:
	@ BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose run truffle bash -c 'echo -e "#run this:\ntruffle deploy --network production --reset" && bash'

run-ui-ganache: ui truffle
	# run a container for ganache-cli
	- docker run --rm -d --name $(GANACHE_CLI) -P \
		-e PRODUCTION_URL -it \
		$(BASEREPO)/$(TRUFFLEIMG) \
		node node_modules/ganache-cli/cli.js --host 0.0.0.0 -d -g 0

	# deploy contract into ganache
	- docker run --rm \
		-it \
		--link $(GANACHE_CLI) \
		$(BASEREPO)/$(TRUFFLEIMG) \
		truffle deploy --network ganachetest --reset

	- docker run -it \
		-e 'CLIENT_URL=ws://ganache-test-incident-reporting:8545' \
		-e 'CLIENT_CONTRACT_ADDRESS=0xcfeb869f69431e42cdb54a4f4f105c19c080a601' \
		--link $(GANACHE_CLI) \
		-p 8080:80 \
		$(BASEREPO)/$(UIIMG)

	- docker stop $(GANACHE_CLI)
	- docker rm $(GANACHE_CLI)

run-ui:
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose up ui

deploy-contract:
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose up -d truffle
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose up deploy-contract

stop:
	BASEREPO=$(BASEREPO) TAG=$(TAG) docker-compose down