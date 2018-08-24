package graph

import (
	context "context"
	"errors"
	"log"
	"time"
)

var (
	ErrInvalidParameter = errors.New("Invalid parameter")
)

type mutationResolver struct{ *GraphQLServer }

func (server *GraphQLServer) Mutation() MutationResolver {
	return &mutationResolver{server}
}

func (server *GraphQLServer) Register(ctx context.Context, args RegisterArgs) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := server.authClient.Register(ctx, args.FirstName, args.LastName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &User{
		ID:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
	}, nil
}
