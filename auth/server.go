//go:generate protoc ./auth.proto --go_out=plugins=grpc:./pb
package auth

import (
	"context"
	"net"

	"github.com/dakhipp/graphql-services/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
}

// ListenGRPC takes a gRPC service and a formatted port string (":8000") and starts up a gRPC server
func ListenGRPC(service Service, port string) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, &grpcServer{service})
	reflection.Register(server)
	return server.Serve(listen)
}

// Register is a gRPC function which registers a user
func (s *grpcServer) Register(ctx context.Context, args *pb.RegisterRequest) (*pb.AuthResponse, error) {
	resp, err := s.service.Register(ctx, args)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{User: &pb.Auth{
		Id:            resp.ID,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		Email:         resp.Email,
		Phone:         resp.Phone,
		Roles:         resp.Roles,
		EmailVerified: resp.EmailVerified,
		PhoneVerified: resp.PhoneVerified,
	}}, nil
}

// Login is a gRPC function that logs in a user
func (s *grpcServer) Login(ctx context.Context, args *pb.LoginRequest) (*pb.AuthResponse, error) {
	resp, err := s.service.Login(ctx, args)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{User: &pb.Auth{
		Id:            resp.ID,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		Email:         resp.Email,
		Phone:         resp.Phone,
		Roles:         resp.Roles,
		EmailVerified: resp.EmailVerified,
		PhoneVerified: resp.PhoneVerified,
	}}, nil
}

// TriggerVerifyEmail is a gRPC function that initiates the email verification process
func (s *grpcServer) TriggerVerifyEmail(ctx context.Context, args *pb.TriggerVerifyEmailRequest) (*pb.MessageResponse, error) {
	if err := s.service.TriggerVerifyEmail(ctx, args); err != nil {
		return nil, err
	}
	return &pb.MessageResponse{
		Message: "A verification link has been sent to your email.",
	}, nil
}

// TriggerVerifyEmail is a gRPC function that initiates the email verification process
func (s *grpcServer) TriggerVerifyPhone(ctx context.Context, args *pb.TriggerVerifyPhoneRequest) (*pb.MessageResponse, error) {
	if err := s.service.TriggerVerifyPhone(ctx, args); err != nil {
		return nil, err
	}
	return &pb.MessageResponse{
		Message: "A verification code has been texted to your phone number.",
	}, nil
}

// VerifyEmail is a gRPC function that completes the email verification process
func (s *grpcServer) VerifyEmail(ctx context.Context, args *pb.VerifyRequest) (*pb.MessageResponse, error) {
	if err := s.service.VerifyEmail(ctx, args); err != nil {
		return nil, err
	}
	return &pb.MessageResponse{
		Message: "Your email has been successfully verified.",
	}, nil
}

// VerifyEmail is a gRPC function that completes the phone verification process
func (s *grpcServer) VerifyPhone(ctx context.Context, args *pb.VerifyRequest) (*pb.MessageResponse, error) {
	if err := s.service.VerifyPhone(ctx, args); err != nil {
		return nil, err
	}
	return &pb.MessageResponse{
		Message: "Your phone has been successfully verified.",
	}, nil
}

// GetUsers is a gRPC function that fetches all users from the database
func (s *grpcServer) GetUsers(ctx context.Context, args *pb.EmptyRequest) (*pb.GetUsersResponse, error) {
	resp, err := s.service.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := []*pb.User{}
	for _, u := range resp {
		users = append(
			users,
			&pb.User{
				Id:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			},
		)
	}
	return &pb.GetUsersResponse{Users: users}, nil
}
