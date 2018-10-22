PROJECT := graphql-services
ECR_PASS := $(shell aws ecr --region=us-west-2 get-authorization-token --output text --query authorizationData[].authorizationToken | base64 --decode | cut -d: -f2)

.CLEAR=\x1b[0m
.BOLD=\x1b[01m

ENV ?= dev
GQL_ENDPOINT ?= http://localhost:8000/graphql

all: help

help:
	@echo "---"
	@echo "$(.BOLD)'$(PROJECT)' Avilable Commands:$(.CLEAR)"
	@echo "---"
	@echo "$(.BOLD)make start$(.CLEAR) - Starts all containers with Docker Compose, you can overwrite ENV environment variable "\
	"to choose which docker-compose.ENV.yaml file to use. \n\tDefault 'ENV=dev' auto-reloads when editing files. \n\tOverwrite to 'ENV=prod' to start containers with compiled binaries."
	@echo "$(.BOLD)make update$(.CLEAR) - Ensures all vendor dependencies are included in the project."
	@echo "$(.BOLD)make gen-all$(.CLEAR) - Re-generates all generated code."
	@echo "$(.BOLD)make gen-auth$(.CLEAR) - Re-generates the auth service's protobuf file."
	@echo "$(.BOLD)make gen-gql$(.CLEAR) - Re-generates the graphql services schema.graphql file using gqlgen."
	@echo "$(.BOLD)make test$(.CLEAR) - Runs integration and unit tests."
	@echo "$(.BOLD)make test-unit$(.CLEAR) - Runs integration tests."
	@echo "$(.BOLD)make test-int$(.CLEAR) - Runs unit tests."
	@echo "$(.BOLD)make tf-init$(.CLEAR) - Allows for '$(.BOLD)terraform init$(.CLEAR)' to be called from the root directory."
	@echo "$(.BOLD)make tf-apply$(.CLEAR) - Allows for '$(.BOLD)terraform apply$(.CLEAR)' to be called from the root directory."
	@echo "$(.BOLD)make tf-destory$(.CLEAR) - Allows for '$(.BOLD)terraform destory$(.CLEAR)' to be called from the root directory."
	@echo "$(.BOLD)make ecr-auth$(.CLEAR) - Authenticates with ERC."
	@echo "$(.BOLD)make ecr-create-all$(.CLEAR) - Authenticates with ECR and creates all required image repositories."
	@echo "$(.BOLD)make ecr-deploy-all$(.CLEAR) - Authenticates with ECR. Builds and deploys all images to ECR."
	@echo "$(.BOLD)make ecr-gql-deploy$(.CLEAR) - Builds and deploys graphql image to ECR."
	@echo "$(.BOLD)make ecr-auth-deploy$(.CLEAR) - Builds and deploys auth image to ECR."
	@echo "$(.BOLD)make ecr-email-deploy$(.CLEAR) - Builds and deploys email image to ECR."
	@echo "$(.BOLD)make ecr-text-deploy$(.CLEAR) - Builds and deploys text image to ECR."
	@echo "$(.BOLD)make kube-set-ecr-secret$(.CLEAR) - Sets ECR authentication key as a Kubernetes secret so it can be used when pulling images from ECR."
	@echo "$(.BOLD)make kube-apply$(.CLEAR) - Applies Kubernetes deployment from ./deployments/k8s/"
	@echo "$(.BOLD)make kube-delete$(.CLEAR) - Deletes Kubernetes deployment from ./deployments/k8s/"
	@echo "$(.BOLD)make minikube-new$(.CLEAR) - Starts a local minikube cluster, runs required commands, and applies Kubernetes deployment from ./deployments/k8s/"
	@echo "$(.BOLD)make minikube-start$(.CLEAR) - Starts a local minikube cluster."
	@echo "$(.BOLD)make minikube-stop$(.CLEAR) - Stops a local minikube cluster."
	@echo "$(.BOLD)make minikube-delete$(.CLEAR) - Deletes a local minikube cluster."
	@echo "$(.BOLD)make minikube-ingress-init$(.CLEAR) - Creates nginx ingress and enables the addon for local minikube cluster so that resources are reachable on host machine."
	@echo "$(.BOLD)make minikube-dash$(.CLEAR) - Creates and launches the Kubernetes dashboard on minikube cluster."

start:
	@docker-compose -f deployments/docker-compose/docker-compose.$(ENV).yml up --build

update:
	@dep ensure

build:
	@make build-auth ; make build-gql

build-auth:
	@go generate auth/server.go

build-gql:
	@go generate graphql/graph/graph.go

test:
	@make test-unit ; make test-int

test-unit:
	@echo "unit tests not implemented yet"

test-int:
	@cd integration ; GQL_ENDPOINT=$(GQL_ENDPOINT) go test ; cd ..

tf-init:
	@cd deployments/terraform/aws/stages/dev ; terraform init ; cd ../../../../..

tf-apply:
	@cd deployments/terraform/aws/stages/dev ; terraform apply ; cd ../../../../..

tf-destroy:
	@cd deployments/terraform/aws/stages/dev ; terraform destroy ; cd ../../../../..

ecr-auth:
	@aws ecr get-login --no-include-email --region us-west-2 | sh

ecr-create-all:
	@make ecr-auth
	@aws ecr create-repository --repository-name dakhipp/graphql
	@aws ecr create-repository --repository-name dakhipp/auth
	@aws ecr create-repository --repository-name dakhipp/email
	@aws ecr create-repository --repository-name dakhipp/text

ecr-deploy-all:
	@make ecr-auth ; make ecr-gql-deploy ; make ecr-auth-deploy ; make ecr-email-deploy ; make ecr-text-deploy

ecr-gql-deploy:
	@docker build -f deployments/docker/graphql/prod.Dockerfile -t dakhipp/graphql .
	@docker tag dakhipp/graphql:latest 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/graphql:latest
	@docker push 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/graphql:latest

ecr-auth-deploy:
	@docker build -f deployments/docker/auth/prod.Dockerfile -t dakhipp/auth .
	@docker tag dakhipp/auth:latest 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/auth:latest
	@docker push 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/auth:latest

ecr-email-deploy:
	@docker build -f deployments/docker/email/prod.Dockerfile -t dakhipp/email .
	@docker tag dakhipp/email:latest 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/email:latest
	@docker push 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/email:latest

ecr-text-deploy:
	@docker build -f deployments/docker/text/prod.Dockerfile -t dakhipp/text .
	@docker tag dakhipp/text:latest 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/text:latest
	@docker push 690303654955.dkr.ecr.us-west-2.amazonaws.com/dakhipp/text:latest

kube-set-ecr-secret:
	@kubectl delete secret --ignore-not-found us-west-2-ecr-registry
	@kubectl create secret docker-registry us-west-2-ecr-registry \
		--docker-server=https://690303654955.dkr.ecr.us-west-2.amazonaws.com \
		--docker-username=AWS \
		--docker-password=$(ECR_PASS) \
		--docker-email="dakhipp@gmail.com"

kube-apply:
	@kubectl apply -f ./deployments/k8s/

kube-delete:
	@kubectl delete -f ./deployments/k8s/

minikube-new:
	@make minikube-start
	@make minikube-ingress-init
	@make kube-set-ecr-secret
	@make kube-apply
	@make minikube-dash

minikube-start:
	@minikube start

minikube-stop:
	@minikube stop

minikube-delete:
	@minikube delete

minikube-ingress-init:
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
	@minikube addons enable ingress

minikube-dash:
	@minikube dashboard
