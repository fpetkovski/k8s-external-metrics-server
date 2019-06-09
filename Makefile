.ONESHELL:
.PHONY: build deploy delete

VERSION?=latest

build:
	eval $$(minikube docker-env)
	docker build -f build/Dockerfile -t external-metrics-server .

deploy:
	kubectl apply -f build/external-metrics-server.yaml
	kubectl apply -f build/beanstalkd.yaml
	kubectl apply -f build/hpa.yaml

delete:
	kubectl delete -f build/beanstalkd.yaml
	kubectl delete -f build/external-metrics-server.yaml
	kubectl delete -f build/hpa.yaml
	eval $$(minikube docker-env)
	docker rmi external-metrics-server