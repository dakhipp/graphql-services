package graph

import (
	context "context"
	"errors"
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth"
)

var (
	ErrInvalidParameter = errors.New("Invalid parameter")
)

func (s *GraphQLServer) Mutation_register(ctx context.Context, args RegisterArgs) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := s.authClient.Register(args)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &User{
		Id:        a.Id,
		FirstName: a.FirstName,
		LastName:  a.LastName,
	}, nil
}
