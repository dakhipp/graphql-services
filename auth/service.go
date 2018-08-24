package auth

import (
	"context"

	"github.com/segmentio/ksuid"
)

// Service : Business logic layer
type Service interface {
	Register(ctx context.Context, firstName string, lastName string) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}

// User : User model
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type authService struct {
	repository Repository
}

// NewService : Creates a new auth services with a repository responsible for database calls
func NewService(repository Repository) Service {
	return &authService{repository}
}

// Register : Creates UUID, attempts to save user to database and returns the user if succesful
func (service *authService) Register(ctx context.Context, firstName string, lastName string) (*User, error) {
	user := &User{
		FirstName: firstName,
		LastName:  lastName,
		ID:        ksuid.New().String(),
	}
	if err := service.repository.CreateUser(ctx, *user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers : Calls repository function to get users from the database
func (service *authService) GetUsers(ctx context.Context) ([]User, error) {
	return service.repository.ReadUsers(ctx)
}
