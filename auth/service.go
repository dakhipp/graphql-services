package auth

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	TriggerVerifyEmail(ctx context.Context, args *pb.TriggerVerifyEmailRequest) error
	TriggerVerifyPhone(ctx context.Context, args *pb.TriggerVerifyPhoneRequest) error
	VerifyEmail(ctx context.Context, args *pb.VerifyRequest) error
	VerifyPhone(ctx context.Context, args *pb.VerifyRequest) error
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
func (s *authService) Register(ctx context.Context, args *pb.RegisterRequest) (*User, error) {
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
	if err := s.repository.CreateUser(ctx, *user); err != nil {
		return nil, err
	}

	a := &pb.TriggerVerifyEmailRequest{
		Email:     args.Email,
		FirstName: args.FirstName,
	}
	if err := s.TriggerVerifyEmail(ctx, a); err != nil {
		return nil, err
	}

	return user, nil
}

// Login fetches a user from the database by email and compares the password they provided with the fetched user's
func (s *authService) Login(ctx context.Context, args *pb.LoginRequest) (*User, error) {
	u, err := s.repository.GetUserByEmail(ctx, args.Email)
	if err != nil {
		return nil, err
	}
	if !checkPasswordHash(args.Password, u.Password) {
		return nil, fmt.Errorf("Invalid credentials")
	}
	return u, nil
}

// TriggerVerifyEmail creates and emailVerification typed code in the database and then produces a message to Kafka in order to send out an email with the code in it
func (s *authService) TriggerVerifyEmail(ctx context.Context, args *pb.TriggerVerifyEmailRequest) error {
	code := generateCode("emailVerification")
	if err := s.repository.CreateVerificationCode(ctx, code); err != nil {
		return err
	}
	if err := s.kafkaProducer.RegisterEmail(ctx, args, code.Code); err != nil {
		return err
	}
	return nil
}

// TriggerVerifyPhone creates and phoneVerification typed code in the database and then produces a message to Kafka in order to send out an text with the code in it
func (s *authService) TriggerVerifyPhone(ctx context.Context, args *pb.TriggerVerifyPhoneRequest) error {
	code := generateCode("phoneVerification")
	if err := s.repository.CreateVerificationCode(ctx, code); err != nil {
		return err
	}
	if err := s.kafkaProducer.ConfirmPhone(ctx, args, code.Code); err != nil {
		return err
	}
	return nil
}

// VerifyEmail checks an email verification code and then updates the users emailVerified property if the code is valid
func (s *authService) VerifyEmail(ctx context.Context, args *pb.VerifyRequest) error {
	if err := s.repository.CheckEmailVerificationCode(ctx, args.Code); err != nil {
		return err
	}
	if err := s.repository.UpdateEmailVerified(ctx, args.UserId, true); err != nil {
		return err
	}
	return nil
}

// VerifyPhone checks an phone verification code and then updates the users phoneVerified property if the code is valid
func (s *authService) VerifyPhone(ctx context.Context, args *pb.VerifyRequest) error {
	if err := s.repository.CheckPhoneVerificationCode(ctx, args.Code); err != nil {
		return err
	}
	if err := s.repository.UpdatePhoneVerified(ctx, args.UserId, true); err != nil {
		return err
	}
	return nil
}

// GetUsers calls a repository function to get users from the database
func (s *authService) GetUsers(ctx context.Context) ([]User, error) {
	return s.repository.ReadUsers(ctx)
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

// randSixDigit returns a random 6 digit string
func randSixDigit() string {
	rand.Seed(time.Now().UnixNano())
	opts := []rune("0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = opts[rand.Intn(len(opts))]
	}
	return string(b)
}

// generate code takes string t which identifies the type of code to create and returns the code
func generateCode(t string) Code {
	return Code{
		Code:    randSixDigit(),
		Type:    t,
		Created: time.Now(),
	}
}
