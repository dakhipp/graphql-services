package graph

import (
	context "context"
	"log"
	"time"
)

func (server *GraphQLServer) Query_getUsers(ctx context.Context) ([]User, error) {
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
