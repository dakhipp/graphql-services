PROJECT := graphql-services

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
	@echo "$(.BOLD)make update$(.CLEAR) - Runs '$(.BOLD)dep ensure$(.CLEAR)' in order to ensure that all required vendor dependencies are included in the project."
	@echo "$(.BOLD)make build$(.CLEAR) - Runs '$(.BOLD)make build-auth$(.CLEAR)' and '$(.BOLD)make build-graphql$(.CLEAR)'"
	@echo "$(.BOLD)make build-auth$(.CLEAR) - Recompiles the auth services gRPC protobuf definition file."
	@echo "$(.BOLD)make build-gql$(.CLEAR) - Recompiles the graphql services schema.graphql file using gqlgen."
	@echo "$(.BOLD)make test$(.CLEAR) - Runs '$(.BOLD)make test-unit$(.CLEAR)' and '$(.BOLD)make test-int$(.CLEAR)'"
	@echo "$(.BOLD)make test-unit$(.CLEAR) - Runs integration tests."
	@echo "$(.BOLD)make test-int$(.CLEAR) - Runs unit tests."
	@echo "$(.BOLD)make tf-init$(.CLEAR) - Convenience command for '$(.BOLD)terraform init$(.CLEAR)'. So it can easily be called from the root directory."
	@echo "$(.BOLD)make tf-apply$(.CLEAR) - Convenience command for '$(.BOLD)terraform apply$(.CLEAR)'. So it can easily be called from the root directory."
	@echo "$(.BOLD)make tf-destory$(.CLEAR) - Convenience command for '$(.BOLD)terraform destory$(.CLEAR)'. So it can easily be called from the root directory."

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
