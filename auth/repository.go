package auth

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	collection = "users"
	database   = "graphql-services"
)

// User model
type User struct {
	ID            string   `bson:"_id" json:"id"`
	FirstName     string   `bson:"firstName" json:"firstName"`
	LastName      string   `bson:"lastName" json:"lastName"`
	Email         string   `bson:"email" json:"email"`
	Phone         string   `bson:"phone" json:"phone"`
	Password      string   `bson:"password" json:"password"`
	Roles         []string `bson:"roles" json:"roles"`
	EmailVerified bool     `bson:"emailVerified" json:"emailVerified"`
	PhoneVerified bool     `bson:"phoneVerified" phoneVerified"`
}

// Mongo is an interface which allows interactions with MongoDB
type Mongo interface {
	CreateUser(ctx context.Context, args User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ReadUsers(ctx context.Context) ([]User, error)
}

type mongoRepository struct {
	db *mgo.Database
}

// NewMongoDBRepository initializes a new database connection for the repository
func NewMongoDBRepository(url string) (Mongo, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	db := session.DB(database)
	c := db.C(collection)

	err = createCollectionIndexes(c)
	if err != nil {
		return nil, err
	}

	return &mongoRepository{db}, err
}

func createCollectionIndexes(c *mgo.Collection) error {
	// create indexes for each indexed field
	for _, key := range []string{"email", "phone"} {
		index := mgo.Index{
			Key:    []string{key},
			Unique: true,
		}
		if err := c.EnsureIndex(index); err != nil {
			return err
		}
	}
	return nil
}

// CreateUser creates a user in the database
func (r *mongoRepository) CreateUser(ctx context.Context, args User) error {
	err := r.db.C(collection).Insert(&args)
	return err
}

func (r *mongoRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.C(collection).Find(bson.M{"email": email}).One(&user)
	return &user, err
}

// // ReadUsers reads users from the database
func (r *mongoRepository) ReadUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.C(collection).Find(bson.M{}).All(&users)
	return users, err
}
