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

// ListenGRPC : Takes a service and a formatted port string (":8000") and starts up a GRPC server
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

// Register : GRPC function which registers a user
func (server *grpcServer) Register(ctx context.Context, args *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	resp, err := server.service.Register(ctx, args)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{User: &pb.Auth{
		Id:    resp.ID,
		Roles: resp.Roles,
	}}, nil
}

// Login is a GRPC function that logs in a user
func (server *grpcServer) Login(ctx context.Context, args *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := server.service.Login(ctx, args)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{User: &pb.Auth{
		Id:    resp.ID,
		Roles: resp.Roles,
	}}, nil
}

// Register : GRPC function which fetches users
func (server *grpcServer) GetUsers(ctx context.Context, args *pb.EmptyRequest) (*pb.GetUsersResponse, error) {
	resp, err := server.service.GetUsers(ctx)
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
