package auth

import (
	"context"

	"github.com/dakhipp/graphql-services/auth/pb"
	"google.golang.org/grpc"
)

// Client : grpc client used to make authentication service calls
type Client struct {
	conn    *grpc.ClientConn
	service pb.AuthServiceClient
}

// NewClient : creates a new grpc client
func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &Client{conn, client}, nil
}

// Close : closes a grpc client connection
func (client *Client) Close() {
	client.conn.Close()
}

// Register : Register function available on the grpc client, registers a user in the database
func (client *Client) Register(ctx context.Context, r *pb.RegisterRequest) (*User, error) {
	resp, err := client.service.Register(
		ctx,
		r,
	)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:    resp.User.Id,
		Roles: resp.User.Roles,
	}, nil
}

// Login is a function available on the GRPC client, it logs in a user
func (client *Client) Login(ctx context.Context, r *pb.LoginRequest) (*User, error) {
	resp, err := client.service.Login(
		ctx,
		r,
	)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:    resp.User.Id,
		Roles: resp.User.Roles,
	}, nil
}

// GetUsers : GetUsers function avilable on the grpc client, fetches all users from the database
func (client *Client) GetUsers(ctx context.Context) ([]User, error) {
	resp, err := client.service.GetUsers(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	users := []User{}
	for _, u := range resp.Users {
		users = append(users, User{
			ID:        u.Id,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		})
	}
	return users, nil
}
