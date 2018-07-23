//go:generate gqlgen -schema ../schema.graphql
package graph

import (
	"github.com/dakhipp/graphql-services/auth"
)

type GraphQLServer struct {
	authClient *auth.Client
}

func NewGraphQLServer(authURL string) (*GraphQLServer, error) {
	// Connect to auth service
	authClient, err := auth.NewClient(authUrl)
	if err != nil {
		return nil, err
	}

	return &GraphQLServer{
		authClient,
	}, nil
}
