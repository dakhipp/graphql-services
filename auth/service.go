package auth

import (
	"context"
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth/pb"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/ksuid"
	"github.com/tinrab/retry"
)

// Service : Business logic layer
type Service interface {
	Register(ctx context.Context, args *pb.RegisterRequest) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}

// configuration struct created from environment variables
type envConfig struct {
	PSQLURL string `envconfig:"PSQL_URL"`
}

// User : User model
type User struct {
	ID        string   `json:"id"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
}

type authService struct {
	repository Repository
}

// New creates and returns a new authService
func New() Service {
	// Declare and attempt to cast config struct
	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Attempt to create auth repository
	var repository Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = NewPostgresRepository(cfg.PSQLURL)
		if err != nil {
			log.Println(err)
		}
		return
	})

	return &authService{repository}
}

// Register : Creates UUID, attempts to save user to database and returns the user if successful
func (service *authService) Register(ctx context.Context, args *pb.RegisterRequest) (*User, error) {
	user := &User{
		ID:        ksuid.New().String(),
		FirstName: args.FirstName,
		LastName:  args.LastName,
		Email:     args.Email,
		Phone:     args.Phone,
		Password:  args.Password,
		Roles:     []string{"USER"},
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
