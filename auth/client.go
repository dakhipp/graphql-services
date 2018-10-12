package auth

import (
	"context"

	"github.com/dakhipp/graphql-services/auth/pb"
	"google.golang.org/grpc"
)

// Client is a struct containing a gRPC connection and a auth service
type Client struct {
	conn    *grpc.ClientConn
	service pb.AuthServiceClient
}

// NewClient creates a new gRPC client
func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &Client{conn, client}, nil
}

// Close closes a gRPC client connection
func (c *Client) Close() {
	c.conn.Close()
}

// Register is a function available on the gRPC client that registers a user in the database
func (c *Client) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.AuthResponse, error) {
	resp, err := c.service.Register(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Login is a function available on the gRPC client that logs in a user
func (c *Client) Login(ctx context.Context, r *pb.LoginRequest) (*pb.AuthResponse, error) {
	resp, err := c.service.Login(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TriggerVerifyEmail is a function available on the gRPC client that initiates the email verification process
func (c *Client) TriggerVerifyEmail(ctx context.Context, email string) (*pb.MessageResponse, error) {
	resp, err := c.service.TriggerVerifyEmail(ctx, &pb.TriggerVerifyEmailRequest{Email: email})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TriggerVerifyPhone is a function available on the gRPC client that initiates the phone verification process
func (c *Client) TriggerVerifyPhone(ctx context.Context, phone string) (*pb.MessageResponse, error) {
	resp, err := c.service.TriggerVerifyPhone(ctx, &pb.TriggerVerifyPhoneRequest{Phone: phone})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// VerifyEmail is a function available on the gRPC client that completes the email verification process
func (c *Client) VerifyEmail(ctx context.Context, args *pb.VerifyRequest) (*pb.MessageResponse, error) {
	resp, err := c.service.VerifyEmail(ctx, args)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// VerifyPhone is a function available on the gRPC client that completes the phone verification process
func (c *Client) VerifyPhone(ctx context.Context, args *pb.VerifyRequest) (*pb.MessageResponse, error) {
	resp, err := c.service.VerifyPhone(ctx, args)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUsers is a function available on the gRPC client that fetches all users from the database
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	resp, err := c.service.GetUsers(ctx, &pb.EmptyRequest{})
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
