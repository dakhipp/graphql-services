//go:generate protoc ./auth.proto --go_out=plugins=grpc:./pb
package auth

import (
	"context"
	"fmt"
	"net"

	"github.com/dakhipp/graphql-services/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAuthServiceServer(serv, &grpcServer{s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	a, err := s.service.Register(ctx, r.FirstName, r.LastName)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{User: &pb.User{
		Id:        a.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
	}}, nil
}

func (s *grpcServer) GetUsers(ctx context.Context) (*pb.GetUsersResponse, error) {
	res, err := s.service.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := []*pb.User{}
	for _, p := range res {
		users = append(
			users,
			&pb.User{
				Id:        p.ID,
				FirstName: p.FirstName,
				LastName:  p.LastName,
			},
		)
	}
	return &pb.GetUsersResponse{Users: users}, nil
}
