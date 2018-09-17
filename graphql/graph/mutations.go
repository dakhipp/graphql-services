package graph

import (
	context "context"
	"fmt"
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth/pb"
	"github.com/segmentio/ksuid"
)

type mutationResolver struct{ *GraphQLServer }

// Mutation : Mutation root resolver, required to satisfy interface
func (server *GraphQLServer) Mutation() MutationResolver {
	return &mutationResolver{server}
}

// Register : Register mutation exposed via GraphQL
func (server *GraphQLServer) Register(ctx context.Context, args RegisterArgs) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// user is attached to context, remove this
	ss, _ := ctx.Value(CONTEXT_SESSION_KEY).(Session)
	fmt.Println(ss)

	r := &pb.RegisterRequest{
		FirstName:    args.FirstName,
		LastName:     args.LastName,
		Email:        args.Email,
		Phone:        args.Phone,
		Password:     args.Password,
		PasswordConf: args.PasswordConf,
	}

	resp, err := server.authClient.Register(ctx, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s := &Session{
		ID:    resp.ID,
		Roles: toRoles(resp.Roles),
	}

	sID := ksuid.New().String()
	server.redisRepository.CreateSession(sID, s)

	server.writeSessionCookie(ctx, sID)

	return s, nil
}

func toRoles(s []string) []Roles {
	c := make([]Roles, len(s))
	for i, v := range s {
		c[i] = Roles(v)
	}
	return c
}
