package graph

import (
	context "context"
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth"
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

	return server.createSession(ctx, resp), nil
}

// Login is a mutation resolver
func (server *GraphQLServer) Login(ctx context.Context, args LoginArgs) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	r := &pb.LoginRequest{
		Email:    args.Email,
		Password: args.Password,
	}

	resp, err := server.authClient.Login(ctx, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return server.createSession(ctx, resp), nil
}

func (server *GraphQLServer) createSession(ctx context.Context, resp *auth.User) *Session {
	s := &Session{
		ID:    resp.ID,
		Roles: toRoles(resp.Roles),
	}

	sID := ksuid.New().String()
	server.redisRepository.CreateSession(sID, s)

	server.writeSessionCookie(ctx, sID)

	return s
}

func toRoles(s []string) []Roles {
	c := make([]Roles, len(s))
	for i, v := range s {
		c[i] = Roles(v)
	}
	return c
}
