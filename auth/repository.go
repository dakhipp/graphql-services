package auth

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq" // import postgres driver
)

// Repository : A repository object which allows interactions with the database
type Repository interface {
	Close()
	CreateUser(ctx context.Context, args User) error
	ReadUsers(ctx context.Context) ([]User, error)
}

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository : Initiates a new database connection for the repository
func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

// Close : Closes the connection for the repository database
func (repository *postgresRepository) Close() {
	repository.db.Close()
}

// Ping : Checks for an error when pinging the database
func (repository *postgresRepository) Ping() error {
	return repository.db.Ping()
}

// CreateUser : Creates a user in the database
func (repository *postgresRepository) CreateUser(ctx context.Context, args User) error {
	_, err := repository.db.ExecContext(ctx, "INSERT INTO users(id, firstName, lastName) VALUES($1, $2, $3)", args.ID, args.FirstName, args.LastName)
	return err
}

// ReadUsers : Reads users from the database
func (repository *postgresRepository) ReadUsers(ctx context.Context) ([]User, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		"SELECT * FROM users ORDER BY id DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u := &User{}
		if err = rows.Scan(&u.ID, &u.FirstName, &u.LastName); err == nil {
			users = append(users, *u)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
