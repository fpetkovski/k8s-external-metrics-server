.ONESHELL:

VERSION?=latest

build:
	eval $$(minikube docker-env)
	docker build -f k8s/Dockerfile -t external-metrics-server src/

deploy:
	kubectl apply -f k8s/external-metrics-server.yaml
	kubectl apply -f k8s/beanstalkd.yaml
	kubectl apply -f k8s/hpa.yaml

delete:
	kubectl delete -f k8s/beanstalkd.yaml
	kubectl delete -f k8s/external-metrics-server.yaml
	kubectl delete -f k8s/hpa.yaml
	eval $$(minikube docker-env)
	docker rmi external-metrics-server