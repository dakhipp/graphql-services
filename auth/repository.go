package auth

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateUser(ctx context.Context, a User) error
	ReadUsers(ctx context.Context) ([]User, error)
}

type postgresRepository struct {
	db *sql.DB
}

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

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) CreateUser(ctx context.Context, a User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users(id, firstName, lastName) VALUES($1, $2, $3)", a.ID, a.FirstName, a.LastName)
	return err
}

func (r *postgresRepository) ReadUsers(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT * FROM users ORDER BY id DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		a := &User{}
		if err = rows.Scan(&a.ID, &a.FirstName, &a.LastName); err == nil {
			users = append(users, *a)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
