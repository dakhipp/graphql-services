//go:generate gorunpkg github.com/99designs/gqlgen --schema ../schema.graphql
package graph

import (
	"github.com/dakhipp/graphql-services/auth"
)

// GraphQLServer : Server containing GRPC client connections
type GraphQLServer struct {
	authClient *auth.Client
}

// NewGraphQLServer : Function that creates a new GraphQLServer from URLs passed into it
func NewGraphQLServer(authURL string) (*GraphQLServer, error) {
	// Connect to auth service
	authClient, err := auth.NewClient(authURL)
	if err != nil {
		return nil, err
	}

	return &GraphQLServer{
		authClient,
	}, nil
}
