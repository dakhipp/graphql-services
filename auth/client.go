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
	c := pb.NewAuthServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Register(ctx context.Context, firstName string, lastName string) (*User, error) {
	r, err := c.service.Register(
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
		ID:        r.User.Id,
		FirstName: r.User.FirstName,
		LastName:  r.User.LastName,
	}, nil
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	r, err := c.service.GetUsers(
		ctx,
		&pb.EmptyRequest{},
	)
	if err != nil {
		return nil, err
	}
	users := []User{}
	for _, p := range r.Users {
		users = append(users, User{
			ID:        p.Id,
			FirstName: p.FirstName,
			LastName:  p.LastName,
		})
	}
	return users, nil
}
