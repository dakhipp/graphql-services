package graph

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Redis interface {
	CreateSession(sID string, s *Session) error
	GetSession(sID string) (Session, error)
	DeleteSession(sID string) error
}

type redisRepository struct {
	client redis.Conn
}

// NewRedisRepository initiates a new redis connection
func NewRedisRepository(url string) (Redis, error) {
	client, err := redis.DialURL(url)
	if err != nil {
		panic(err)
	}
	return &redisRepository{client}, nil
}

// CreateSession takes a session ID and a session and creates that session in Redis
func (r *redisRepository) CreateSession(sID string, s *Session) error {
	// marshal session into byte array
	b, _ := json.Marshal(s)
	// expires in 7 days
	e := "604800"
	// save byte array into redis, if there is an error return it
	_, err := r.client.Do("SET", sID, b, "EX", e)
	if err != nil {
		return err
	}
	return nil
}

// GetSession takes a session ID and looks up a session in Redis
func (r *redisRepository) GetSession(sID string) (Session, error) {
	// get byte representation of session from redis
	b, err := redis.Bytes(r.client.Do("GET", sID))
	if err != nil {
		fmt.Println(err)
		return Session{}, err
	}
	// unmarshal bytes into a Session and return it
	var s Session
	json.Unmarshal(b, &s)
	return s, nil
}

// DeleteSession takes a session ID and uses it to delete the session from Redis
func (r *redisRepository) DeleteSession(sID string) error {
	_, err := r.client.Do("DEL", sID)
	return err
}
