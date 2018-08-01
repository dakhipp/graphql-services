package auth

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	Register(ctx context.Context, firstName string, lastName string) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type authService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &authService{repository}
}

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

func (service *authService) GetUsers(ctx context.Context) ([]User, error) {
	return service.repository.ReadUsers(ctx)
}
