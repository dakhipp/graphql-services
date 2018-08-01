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

func (server *GraphQLServer) Mutation_register(ctx context.Context, args RegisterArgs) (*User, error) {
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
