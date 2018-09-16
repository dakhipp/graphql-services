package graph

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Redis interface {
	CreateSession(sId string, s Session) error
	GetSession(sId string) (Session, error)
}

type redisRepository struct {
	client redis.Conn
}

// NewRedisRepository : Initiates a new database connection for the repository
func NewRedisRepository(url string) (Redis, error) {
	client, err := redis.DialURL(url)
	if err != nil {
		panic(err)
	}
	return &redisRepository{client}, nil
}

// CreateSession takes a session ID and a session and creates that session in Redis
func (r *redisRepository) CreateSession(sID string, s Session) error {
	_, err := r.client.Do("HMSET", sID, "ID", s.ID, "Name", s.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetSession takes a session ID and looks up a session in Redis
func (r *redisRepository) GetSession(sID string) (Session, error) {
	result, err := redis.Strings(r.client.Do("HMGET", sID, "ID", "Name"))
	if err != nil {
		fmt.Println(err)
		return Session{}, err
	}
	return Session{
		result[0],
		result[1],
	}, nil
}
