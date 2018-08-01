package auth

import (
	"context"

	"github.com/dakhipp/graphql-services/auth/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AuthServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &Client{conn, client}, nil
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) Register(ctx context.Context, firstName string, lastName string) (*User, error) {
	resp, err := client.service.Register(
		ctx,
		&pb.RegisterRequest{
			FirstName: firstName,
			LastName:  lastName,
		},
	)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        resp.User.Id,
		FirstName: resp.User.FirstName,
		LastName:  resp.User.LastName,
	}, nil
}

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
