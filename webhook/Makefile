IMG ?= webhook-server:1.0.0

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
	cd k8s/webhook-server && kustomize edit set image webhook-server=$(IMG) 
	kustomize build k8s/webhook-server | kubectl apply -f -
	scripts/deployment_status_monitor.sh webhook-server --ns admission-webhook
	kubectl apply -k k8s/webhook-configuration

.PHONY: test
test:
	kubectl apply -f k8s/test

.PHONY: uplift
uplift: image
	kind load docker-image $(IMG) --name demo
	kubectl delete -f k8s/webhook-configuration/webhook-configuration.yaml
	kubectl delete -k k8s/webhook-server
	cd k8s/webhook-server && kustomize edit set image webhook-server=$(IMG) 
	kustomize build k8s/webhook-server | kubectl apply -f -
	scripts/deployment_status_monitor.sh webhook-server --ns admission-webhook
	kubectl apply -k k8s/webhook-configuration

.PHONY: clean
clean:
	rm server/server
	kind delete cluster --name demo