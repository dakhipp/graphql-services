package auth

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	Register(ctx context.Context, firstName string, lastName string) (*User, error)
	GetUsers(ctx context.Context) (*User, error)
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type authService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &authService{r}
}

func (s *authService) Register(ctx context.Context, firstName string, lastName string) (*User, error) {
	a := &User{
		FirstName: firstName,
		LastName:  lastName,
		ID:        ksuid.New().String(),
	}
	if err := s.repository.CreateUser(ctx, *a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *authService) GetUsers(ctx context.Context) ([]User, error) {
	return s.repository.ReadUsers(ctx)
}
