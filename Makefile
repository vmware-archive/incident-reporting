TAG ?= latest
REPO ?= index.docker.io/tompscanlan

all: containers

containers: ui truffle
push: push-ui push-truffle

ui: ui-server/*.go ui-server/Dockerfile truffle/contracts/*.sol
	docker build -t incident-reporting-ui:$(TAG) -f ui-server/Dockerfile .

truffle: truffle/package*.json truffle/Dockerfile truffle/contracts/*.sol truffle/migrations/*.js truffle/test/*.js truffle/*.js
	docker build -t incident-reporting-truffle:$(TAG) -f truffle/Dockerfile .

push-ui: ui
	docker tag incident-reporting-ui:$(TAG) $(REPO)/incident-reporting-ui:$(TAG)
	docker push $(REPO)/incident-reporting-ui:$(TAG)

push-truffle: truffle
	docker tag incident-reporting-truffle:$(TAG) $(REPO)/incident-reporting-truffle:$(TAG)
	docker push $(REPO)/incident-reporting-truffle:$(TAG)
