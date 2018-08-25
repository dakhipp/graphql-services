
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
Each stage will have its own `terraform.tfvars` file to hold credentials. You can use `terraform.example.tfvars` to see what fields are needed. In order to create a bastion server public key you will need to do the following:
1. `ssh-keygen -t rsa -b 4096 -C "<your_email@example.com>" -f $HOME/.ssh/<bastion-key-name>`
2. `cat $HOME/.ssh/<bastion-key-name>.pub | pbcopy`
3. Paste the value in your clipboard into the `bastion_public_key` field of the `terraform.tfvars` file.
