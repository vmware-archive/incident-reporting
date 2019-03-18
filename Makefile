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
	docker build -t $(UIIMG) -f ui-server/Dockerfile .
	docker tag $(UIIMG) $(BASEREPO)/$(UIIMG)

truffle: truffle/package*.json truffle/Dockerfile truffle/contracts/*.sol truffle/migrations/*.js truffle/test/*.js truffle/*.js
	docker build -t $(TRUFFLEIMG) -f truffle/Dockerfile .
	docker tag $(TRUFFLEIMG) $(BASEREPO)/$(TRUFFLEIMG)

push-ui: ui
	docker push $(BASEREPO)/$(UIIMG)

push-truffle: truffle
	docker push $(BASEREPO)/$(TRUFFLEIMG)

run-truffle:
	docker run -e PRODUCTION_URL -it incident-reporting-truffle:$(TAG) bash -c 'echo -e "#run this:\ntruffle deploy --network production --reset" && bash'

run-ui:
	docker run -it \
		--env-file env \
		-p 8080:80 \
		incident-reporting-ui:$(TAG)

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
