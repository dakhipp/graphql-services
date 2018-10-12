package graph

import (
	context "context"
	"fmt"
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth/pb"
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

	// handle validation
	if err := s.validation.validate.Struct(args); err != nil {
		return nil, formatValidationErrors(ctx, err)
	}

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

	// handle writing of the session and then return the session or an error
	return s.handleWriteSession(ctx, resp)
}

// Login is a mutation resolver that logs in and returns a session cookie to identify them with
func (s *GraphQLServer) Login(ctx context.Context, args LoginArgs) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// handle validation
	if err := s.validation.validate.Struct(args); err != nil {
		return nil, formatValidationErrors(ctx, err)
	}

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

	// handle writing of the session and then return the session or an error
	return s.handleWriteSession(ctx, resp)
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
	if err := s.redisRepository.DeleteSession(sID); err != nil {
		fmt.Println(err)
	}

	// return message to the user
	return &Message{
		Message: "You've been logged out.",
	}, nil
}

// TriggerVerifyEmail initiates the email verification process
func (s *GraphQLServer) TriggerVerifyEmail(ctx context.Context) (*Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// send request to the gRPC server
	resp, err := s.authClient.TriggerVerifyEmail(ctx, s.getSessionFromContext(ctx).Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// return message to the user
	return &Message{Message: resp.Message}, nil
}

// TriggerVerifyPhone initiates the email verification process
func (s *GraphQLServer) TriggerVerifyPhone(ctx context.Context) (*Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// send request to the gRPC server
	resp, err := s.authClient.TriggerVerifyPhone(ctx, s.getSessionFromContext(ctx).Phone)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// return message to the user
	return &Message{Message: resp.Message}, nil
}

// VerifyEmail completes the email verification process
func (s *GraphQLServer) VerifyEmail(ctx context.Context, code string) (*Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// cast GraphQL arguments into gRPC VerifyRequest
	r := &pb.VerifyRequest{
		UserId: s.getSessionFromContext(ctx).ID,
		Code:   code,
	}

	// send request to the gRPC server
	resp, err := s.authClient.VerifyEmail(ctx, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// return message to the user
	return &Message{Message: resp.Message}, nil
}

// VerifyPhone completes the email verification process
func (s *GraphQLServer) VerifyPhone(ctx context.Context, code string) (*Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// cast GraphQL arguments into gRPC VerifyRequest
	r := &pb.VerifyRequest{
		UserId: s.getSessionFromContext(ctx).ID,
		Code:   code,
	}

	// send request to the gRPC server
	resp, err := s.authClient.VerifyPhone(ctx, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// return message to the user
	return &Message{Message: resp.Message}, nil
}
