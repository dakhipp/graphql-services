
### System Dependencies:
- Golang v1.10.^
- Dep Golang Package Manager
- Docker & Docker Tool Chain

---

### Commands:
Run `make` in order to see all available commands and their uses.

---

### Service Overview:

##### Publicly Expose Services:
- graphql - A GraphQL API Gateway which handles session storage and authorization.

##### Internal gRPC Exposed Services:
- auth - An authentication and user management service.

##### Internal Un-Exposed Services:
- email - A Kafka consumer that is capable of sending emails when messages are produced to the Kafka topic "email"
- text - A Kafka consumer that is capable of sending text messages when messages are produced to the Kafka topic "text"

---

### Example gRPC Service File Breakdown:
- `./auth/cmd/account/main.go` 	  - Implementation of the auth service. Reads environment variables and starts gRPC server.
- `./auth/pb/account.pb.go` 			- Generated protobuf types and functions used to implement gRPC services.
- `./auth/auth.proto`							- Protobuf definition file used to generate file above.
- `./auth/server.go`							- gRPC server used to expose service functionality.
- `./auth/client.go`							- gRPC client used by the graphql service to make requests to the gRPC server.
- `./auth/repository.go`					- Data store methods used to persist and manipulate data in MongoDB.
- './auth/producer.go'            - Kafka message producer used to produce messages that will be used by services that are not exposed via gRPC.
- `./auth/service.go`							- Business logic layer that handles calling repository functions or producing messages to Kafka.

---

### Service Secrets:
The text and email services require secrets to be loaded in via environment variables. In development this can be done with a .env.dev file which mirrors .env.example.

---

### Terraform Modules:
- `Bastion`       - Bastion server that allows access to resources in public and private subnets
- `CodePipeline`  - CI / CD Pipeline with CodePipeline, CodeBuild, and CodeDeploy
- `ECS`           - Container orchestration, load balancing, log aggregation, and metrics through ECS
- `RDS`           - RDS service with PostgreSQL
- `Route53`       - DNS services with Route53
- `VPC`           - VPC for private networking

### Terraform Secrets:
Each stage will have its own `terraform.tfvars` file to hold credentials. You can use `terraform.example.tfvars` to see what fields are needed. 

In order to create a bastion server public key you will need to do the following:
1. `ssh-keygen -t rsa -b 4096 -C "<your_email@example.com>" -f $HOME/.ssh/<bastion-key-name>`
2. `cat $HOME/.ssh/<bastion-key-name>.pub | pbcopy`
3. Paste the value in your clipboard into the `bastion_public_key` field of the `terraform.tfvars` file.

In order for CodePipeline to automatically pull from a Github branch you will need to create a Github "Personal access token": https://github.com/settings/tokens

You'll want to create a Route53 hosted zone for a domain name and request an SSL certificate from Amazon's ACM service on the domain to host at.

### Terraform Info:
- The VPC, RDS, and Bastion modules can easily be pulled into other projects, however ECS and Codepipeline contain quite a bit of project specific configuration. 
