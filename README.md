
Commands:

Generate protocol file from server file with comment at top:
- `go generate auth/server.go` (comment: //go:generate protoc ./auth.proto --go_out=plugins=grpc:./pb) 

Generate GraphQL files from file with comment at top: 
- `go generate ./graphql/graph/graph.go` (comment: //go:generate gqlgen -schema ../schema.graphql -typemap ../types.json) 


Update dependency: 
- `go get -u github.com/dakhipp...` 


---

Each service has 3 layers:
- server 					(responsible for communication)
- service					(contains business logic)
- repository			(writing and reading data from a database)

Also The Following:
- proto definiton		(.proto definition file)
- client						(grpc client)
- cmd/service 			(grpc server)
- pb/service 				(generated pb file from definition)

---

Skeleton Service:

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

Example:

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

File Breakdown:
- ./account/cmd/account/main.go 	= grpc server
- ./account/pb/account.pb.go 			= generated proto file
- ./account.proto									= protocol definition file used to generate file above
- ./client.go											= grpc client
- ./repository.go									= data access methods
- ./server.go											= http interface for your services
- ./service.go										= business logic layer that accesses your repository functions

---

