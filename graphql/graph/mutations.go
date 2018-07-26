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

func (s *GraphQLServer) Mutation_register(ctx context.Context, args RegisterArgs) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := s.authClient.Register(ctx, args.FirstName, args.LastName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &User{
		ID:        a.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
	}, nil
}
