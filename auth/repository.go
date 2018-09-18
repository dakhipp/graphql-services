package auth

import (
	"context"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Repository : A repository object which allows interactions with the database
type Repository interface {
	Close()
	CreateUser(ctx context.Context, args User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ReadUsers(ctx context.Context) ([]User, error)
}

type postgresRepository struct {
	db *pg.DB
}

// NewPostgresRepository initializes a new database connection for the repository
func NewPostgresRepository(url string) (Repository, error) {
	opts, err := pg.ParseURL(url)
	if err != nil {
		fmt.Println(err)
	}
	db := pg.Connect(opts)

	return &postgresRepository{db}, nil
}

// Close closes the connection
func (repository *postgresRepository) Close() {
	repository.db.Close()
}

// CreateUser creates a user in the database
func (repository *postgresRepository) CreateUser(ctx context.Context, args User) error {
	err := repository.db.Insert(&args)
	return err
}

func (repository *postgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	_, err := repository.db.QueryOne(&u, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ReadUsers reads users from the database
func (repository *postgresRepository) ReadUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := repository.db.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// create schema migrates
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
