package graph

import (
	context "context"
	"log"
	"time"
)

func (s *GraphQLServer) Query_getUsers(ctx context.Context) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	r, err := s.authClient.GetUsers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	users := []User{}
	for _, a := range r {
		users = append(users,
			User{
				ID:        a.ID,
				FirstName: a.FirstName,
				LastName:  a.LastName,
			},
		)
	}

	return users, nil
}
