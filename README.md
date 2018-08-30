
## Commands:
- `go generate auth/server.go`            - Generate protobuf file.
- `go generate ./graphql/graph/graph.go`  - Generate GraphQL files with gqlgen.
- `dep ensure`                            - Update dependencies.

---

## Service File Breakdown:
- `./auth/cmd/account/main.go` 	  - Connects to DB, starts and exposes GRPC server.
- `./auth/pb/account.pb.go` 			- Generated proto functions to communicate between GRPC client & server.
- `./auth.proto`									- Protobuf definition file used to generate file above.
- `./client.go`									  - GRPC client.
- `./repository.go`								- Data access methods.
- `./server.go`										- GRPC server.
- `./service.go`									- Business logic layer that calls repository functions.

---

## Skeleton Service:

```
<SERVICE>
├── cmd
│   └── <SERVICE>
│       └── main.go
├── pb
│   └── <SERVICE>.pb.go
├── <SERVICE>.proto
├── client.go
├── Dockerfile
├── repository.go
├── server.go
├── service.go`
```

---

## Terraform Info:

#### Modules:
- `Bastion`       - Bastion server that allows access to resources in public and private subnets
- `CodePipeline`  - CI / CD Pipeline with CodePipeline, CodeBuild, and CodeDeploy
- `ECS`           - Container orchestration, load balancing, log aggregation, and metrics through ECS
- `RDS`           - RDS service with PostgreSQL
- `Route53`       - DNS services with Route53
- `VPC`           - VPC for private networking

#### Secrets:
Each stage will have its own `terraform.tfvars` file to hold credentials. You can use `terraform.example.tfvars` to see what fields are needed. 

In order to create a bastion server public key you will need to do the following:
1. `ssh-keygen -t rsa -b 4096 -C "<your_email@example.com>" -f $HOME/.ssh/<bastion-key-name>`
2. `cat $HOME/.ssh/<bastion-key-name>.pub | pbcopy`
3. Paste the value in your clipboard into the `bastion_public_key` field of the `terraform.tfvars` file.

In order for CodePipeline to automatically pull from a Github branch you will need to create a Github "Personal access token": https://github.com/settings/tokens

You'll want to create a Route53 hosted zone for a domain name and request an SSL certificate from Amazon's ACM service on the domain to host at.

#### Extra Info:
- The VPC, RDS, and Bastion modules can easily be pulled into other projects, however ECS and Codepipeline contain quite a bit of project specific configuration. 

#### Commands:
- `yarn terraform:init` - Convenience command, cd's into the terraform AWS directory and runs `terraform init`, then cd's back into main project directory.
- `yarn terraform:apply` - Convenience command, cd's into the terraform AWS directory and runs `terraform apply`, then cd's back into main project directory.
- `yarn terraform:destroy` - Convenience command, cd's into the terraform AWS directory and runs `terraform destroy`, then cd's back into main project directory.

---

## Migrations:

#### Dev:
In development you can trigger migration creation, application, and roll back manually using commands from package.json. You must have the dev environment started up, since the commands will run inside docker containers to take advantage of environment variables.

#### Prod:
In production or any external server environment migrations will automatically be applied when deployed. If you need to rollback a migration you'll have to SSH into the cluster and manually apply the migration rollback using the migration binary file inside the migration container.

#### Commands:
- `yarn migrate:create <MIGRATION_NAME>` - Creates a new migration template file. NOTE: Commands are ran in a docker container and require that the dev environment has been started.
- `yarn migrate:up` - Runs the next batch of migrations. NOTE: Commands are ran in a docker container and require that the dev environment has been started.
- `yarn migrate:down` - Rolls back a batch of migrations. NOTE: Commands are ran in a docker container and require that the dev environment has been started.
