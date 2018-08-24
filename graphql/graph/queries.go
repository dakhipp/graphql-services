package graph

import (
	context "context"
	"log"
	"time"
)

type queryResolver struct{ *GraphQLServer }

// Query : Query root resolver, required to satisfy interface
func (server *GraphQLServer) Query() QueryResolver {
	return &queryResolver{server}
}

// GetUsers : GetUsers query exposed via GraphQL
func (server *GraphQLServer) GetUsers(ctx context.Context) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := server.authClient.GetUsers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	users := []User{}
	for _, u := range resp {
		users = append(users,
			User{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			},
		)
	}

	return users, nil
}
