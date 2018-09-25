package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth/pb"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/ksuid"
	"github.com/tinrab/retry"
	"golang.org/x/crypto/bcrypt"
)

// Service is an interface of business logic functions
type Service interface {
	Register(ctx context.Context, args *pb.RegisterRequest) (*User, error)
	Login(ctx context.Context, args *pb.LoginRequest) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}

// configuration struct created from environment variables
type envConfig struct {
	MongoURL  string `envconfig:"MONGO_URL"`
	KafkaAddr string `envconfig:"KAFKA_ADDR"`
}

type authService struct {
	repository    Mongo
	kafkaProducer Kafka
}

// New creates and returns a new authService
func New() Service {
	// Declare and attempt to cast config struct
	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// attempt to create MongoDB connection
	var repository Mongo
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = NewMongoDBRepository(cfg.MongoURL)
		if err != nil {
			log.Println(err)
		}
		return
	})

	// attempt to open a Kafka connection
	var kafkaProducer Kafka
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		kafkaProducer, err = NewKafkaProducer(cfg.KafkaAddr)
		if err != nil {
			log.Println(err)
		}
		return
	})

	return &authService{
		repository,
		kafkaProducer,
	}
}

// Register hashes a users password, creates them in the database, and sends out the registration email
func (service *authService) Register(ctx context.Context, args *pb.RegisterRequest) (*User, error) {
	hp, err := hashPassword(args.Password)
	if err != nil {
		return nil, err
	}

	// build user from RegistrationRequest
	user := &User{
		ID:        ksuid.New().String(),
		FirstName: args.FirstName,
		LastName:  args.LastName,
		Email:     args.Email,
		Phone:     args.Phone,
		Password:  hp,
		Roles:     []string{"USER"},
	}

	// create user in database
	if err := service.repository.CreateUser(ctx, *user); err != nil {
		return nil, err
	}

	// produce a message to kafka so our email service will send out the registration email
	if err := service.kafkaProducer.RegisterEmail(ctx, *user); err != nil {
		return nil, err
	}

	// produce a message to kafka so our email service will send out the registration email
	if err := service.kafkaProducer.ConfirmPhone(ctx, *user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login fetches a user from the database by email and compares the password they provided with the fetched user's
func (service *authService) Login(ctx context.Context, args *pb.LoginRequest) (*User, error) {
	u, err := service.repository.GetUserByEmail(ctx, args.Email)
	if err != nil {
		return nil, err
	}
	if !checkPasswordHash(args.Password, u.Password) {
		return nil, fmt.Errorf("Invalid credentials")
	}
	return u, nil
}

// GetUsers calls a repository function to get users from the database
func (service *authService) GetUsers(ctx context.Context) ([]User, error) {
	return service.repository.ReadUsers(ctx)
}

// hashPassword takes a plain text password and hashes it
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPasswordHash takes a text password and a stored password hash and returns whether or not they match
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
