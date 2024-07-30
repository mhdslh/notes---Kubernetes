IMG ?= webhook-server:0.1.0

CONTAINER_TOOL ?= docker

.PHONY: build
build:
	cd server; export GOTMPDIR=$(pwd); go mod tidy; go build -o server main.go

.PHONY: image
image: build
	$(CONTAINER_TOOL) build --tag $(IMG) .

.PHONY: push
push: image
	$(CONTAINER_TOOL) push $(IMG)

.PHONY: deploy
deploy: image
	kind create cluster --name demo
	kind load docker-image $(IMG) --name demo
	kubectl apply -f resources/webhook-namespace.yaml
	kubectl apply -f resources/webhook-server.yaml
	scripts/deployment_status_monitor.sh webhook-server --ns admission-webhook
	kubectl apply -f resources/webhook-configuration.yaml

.PHONY: test
test:
	kubectl apply -f resources/test-pod.yaml

.PHONY: uplift
uplift: image
	kind load docker-image $(IMG) --name demo
	kubectl delete -f resources/webhook-configuration.yaml
	kubectl delete -f resources/webhook-server.yaml
	kubectl apply -f resources/webhook-server.yaml
	scripts/deployment_status_monitor.sh webhook-server --ns admission-webhook
	kubectl apply -f resources/webhook-configuration.yaml

.PHONY: clean
clean:
	rm server/server
	kind delete cluster --name demo