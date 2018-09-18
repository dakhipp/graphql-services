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

// Mutation is required to satisfy interface
func (s *GraphQLServer) Mutation() MutationResolver {
	return &mutationResolver{s}
}

// Register is a mutation resolver that registers a user into the database and returns a session cookie to identify them with
func (s *GraphQLServer) Register(ctx context.Context, args RegisterArgs) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// cast GraphQL arguments into gRPC RegisterRequest
	r := &pb.RegisterRequest{
		FirstName:    args.FirstName,
		LastName:     args.LastName,
		Email:        args.Email,
		Phone:        args.Phone,
		Password:     args.Password,
		PasswordConf: args.PasswordConf,
	}

	// send request to the gRPC server
	resp, err := s.authClient.Register(ctx, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// create unique session ID
	sID := ksuid.New().String()

	// create session
	ses := s.createSession(ctx, resp)

	// create session cookie
	c := s.createSessionCookie(ctx, sID)

	// use http.ResponseWriter to write the cookie into the response
	s.writeSessionCookie(ctx, c)

	// save session in redis
	s.redisRepository.CreateSession(sID, ses)

	// return the session
	return ses, nil
}

// Login is a mutation resolver that logs in and returns a session cookie to identify them with
func (s *GraphQLServer) Login(ctx context.Context, args LoginArgs) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// cast GraphQL arguments into gRPC LoginRequest
	r := &pb.LoginRequest{
		Email:    args.Email,
		Password: args.Password,
	}

	// send request to the gRPC server
	resp, err := s.authClient.Login(ctx, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// create unique session ID
	sID := ksuid.New().String()

	// create session
	ses := s.createSession(ctx, resp)

	// create session cookie
	c := s.createSessionCookie(ctx, sID)

	// use http.ResponseWriter to write the cookie into the response
	s.writeSessionCookie(ctx, c)

	// save session in redis
	s.redisRepository.CreateSession(sID, ses)

	// return the session
	return ses, nil
}

// Logout is a mutation resolver that expires a cookie, deletes their session, and returns a success message
func (s *GraphQLServer) Logout(ctx context.Context) (*Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// get the session ID from the request
	sID := s.getSessionID(ctx)

	// expire the session cookie
	c := s.expireSessionCookie(sID)

	// use http.ResponseWriter to write the cookie into the response
	s.writeSessionCookie(ctx, c)

	// delete the session from redis and log if there is an error
	err := s.redisRepository.DeleteSession(sID)
	if err != nil {
		fmt.Println(err)
	}

	// return message to the user
	return &Message{
		Message: "You've been logged out.",
	}, nil
}
