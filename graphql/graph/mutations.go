package graph

import (
	context "context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"
)

type mutationResolver struct{ *GraphQLServer }

// Mutation : Mutation root resolver, required to satisfy interface
func (server *GraphQLServer) Mutation() MutationResolver {
	return &mutationResolver{server}
}

// Register : Register mutation exposed via GraphQL
func (server *GraphQLServer) Register(ctx context.Context, args RegisterArgs) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// user is attached to context, remove this
	ss, _ := ctx.Value(CONTEXT_SESSION_KEY).(Session)
	fmt.Println(ss)

	resp, err := server.authClient.Register(ctx, args.FirstName, args.LastName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	u := &User{
		ID:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
	}

	sID := ksuid.New().String()
	s := server.sessionFromUser(u)
	server.redisRepository.CreateSession(sID, s)

	server.writeSessionCookie(ctx, sID)

	return u, nil
}
