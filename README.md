
## Commands:

Generate protocol file from server file with comment at top:
- `go generate auth/server.go` (comment: //go:generate protoc ./auth.proto --go_out=plugins=grpc:./pb) 

Generate GraphQL files from file with comment at top: 
- `go generate ./graphql/graph/graph.go` (comment: //go:generate gqlgen -schema ../schema.graphql -typemap ../types.json) 


Update dependency: 
- `dep ensure` 


---

## Each service has 3 layers:
- server 					(responsible for communication)
- service					(contains business logic)
- repository			(writing and reading data from a database)

Also The Following:
- proto definition		(.proto definition file)
- client						  (GRPC client)
- cmd/service 			  (GRPC server)
- pb/service 				  (generated pb file from definition)

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

## Example:

```
account
├── cmd
│   └── account
│       └── main.go
├── pb
│   └── account.pb.go
├── account.proto
├── client.go
├── Dockerfile
├── repository.go
├── server.go
├── service.go
```

---

## File Breakdown:
- ./account/cmd/account/main.go 	= GRPC server
- ./account/pb/account.pb.go 			= generated proto file
- ./account.proto									= protocol definition file used to generate file above
- ./client.go											= grpc client
- ./repository.go									= data access methods
- ./server.go											= http interface for your services
- ./service.go										= business logic layer that accesses your repository functions

---

## Terraform Info:

### Secrets:
Each stage will have its own `terraform.tfvars` file to hold credentials. You can use `terraform.example.tfvars` to see what fields are needed. 

In order to create a bastion server public key you will need to do the following:
1. `ssh-keygen -t rsa -b 4096 -C "<your_email@example.com>" -f $HOME/.ssh/<bastion-key-name>`
2. `cat $HOME/.ssh/<bastion-key-name>.pub | pbcopy`
3. Paste the value in your clipboard into the `bastion_public_key` field of the `terraform.tfvars` file.

In order to allow Codepipeline to automatically pull from a Github branch you will need to create a Github `Personal access token`: https://github.com/settings/tokens

You'll want to create a Route53 hosted zone for a domain name and request an SSL certificate from Amazon's ACM service on the domain to host at.

### Extra Info:
- The VPC, RDS, and Bastion modules can easily be pulled into other projects, however ECS and Codepipeline contain quite a bit of project specific configuration. 

### Commands:
- `yarn terraform:init` - Convenience command, cd's into the terraform AWS directory and runs `terraform init`, then cd's back into main project directory.
- `yarn terraform:apply` - Convenience command, cd's into the terraform AWS directory and runs `terraform apply`, then cd's back into main project directory.
- `yarn terraform:destroy` - Convenience command, cd's into the terraform AWS directory and runs `terraform destroy`, then cd's back into main project directory.

---

## Migrations:

### Dev:
In development you can trigger migration creation, application, and roll back manually using commands from package.json. You must have the dev environment started up, since the commands will run inside docker containers to take advantage of environment variables.

### Prod:
In production or any external server environment migrations will automatically be applied when deployed. If you need to rollback a migration you'll have to SSH into the cluster and manually apply the migration rollback using the migration binary file inside the migration container.

### Commands:
- `yarn migrate:create <MIGRATION_NAME>` - Creates a new migration template file. NOTE: Commands are ran in a docker container and require that the dev environment has been started.
- `yarn migrate:up` - Runs the next batch of migrations. NOTE: Commands are ran in a docker container and require that the dev environment has been started.
- `yarn migrate:down` - Rolls back a batch of migrations. NOTE: Commands are ran in a docker container and require that the dev environment has been started.
