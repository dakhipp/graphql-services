package auth

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	userCollection = "users"
	codeCollection = "codes"
	database       = "graphql-services"
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
	PhoneVerified bool     `bson:"phoneVerified" json:"phoneVerified"`
}

// Code model
type Code struct {
	ID      string    `bson:"_id,omitempty" json:"id"`
	Code    string    `bson:"code" json:"code"`
	Type    string    `bson:"type" json:"type"`
	Created time.Time `bson:"created" json:"created"`
}

// Mongo is an interface which allows interactions with MongoDB
type Mongo interface {
	CreateUser(ctx context.Context, args User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ReadUsers(ctx context.Context) ([]User, error)
	CreateVerificationCode(ctx context.Context, code Code) error
	CheckPhoneVerificationCode(ctx context.Context, code string) error
	CheckEmailVerificationCode(ctx context.Context, code string) error
	UpdatePhoneVerified(ctx context.Context, id string, value bool) error
	UpdateEmailVerified(ctx context.Context, id string, value bool) error
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
	uc := db.C(userCollection)
	cc := db.C(codeCollection)

	if err := createUserCollectionIndexes(uc); err != nil {
		return nil, err
	}

	if err := createCodeCollectionIndexes(cc); err != nil {
		return nil, err
	}

	return &mongoRepository{db}, nil
}

// createUserCollectionIndexes creates required indexes on the "users" collection
func createUserCollectionIndexes(c *mgo.Collection) error {
	// create indexes for each indexed field
	for _, key := range []string{"email", "phone"} {
		i := mgo.Index{
			Key:    []string{key},
			Unique: true,
		}
		if err := c.EnsureIndex(i); err != nil {
			return err
		}
	}
	return nil
}

// createCodeCollectionIndexes creates required indexes on the "codes" collection, including the TTL index which deletes documents after 15 minutes
func createCodeCollectionIndexes(c *mgo.Collection) error {
	i := mgo.Index{
		Key:         []string{"created"},
		Unique:      false,
		DropDups:    false,
		Background:  true,
		ExpireAfter: time.Minute * 15,
	}
	if err := c.EnsureIndex(i); err != nil {
		return err
	}
	return nil
}

// CreateUser creates a user in the database
func (r *mongoRepository) CreateUser(ctx context.Context, args User) error {
	err := r.db.C(userCollection).Insert(&args)
	return err
}

// GetUserByEmail fetches a user from the database via email
func (r *mongoRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.C(userCollection).Find(bson.M{"email": email}).One(&user)
	return &user, err
}

// ReadUsers reads users from the database
func (r *mongoRepository) ReadUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.C(userCollection).Find(bson.M{}).All(&users)
	return users, err
}

// CreateVerificationCode creates a new verification code in the database
func (r *mongoRepository) CreateVerificationCode(ctx context.Context, code Code) error {
	err := r.db.C(codeCollection).Insert(&code)
	return err
}

// CheckPhoneVerificationCode checks a verification code with the type "phoneVerification"
func (r *mongoRepository) CheckPhoneVerificationCode(ctx context.Context, code string) error {
	err := r.db.C(codeCollection).Remove(bson.M{"type": "phoneVerification", "code": code})
	return err
}

// CheckEmailVerificationCode checks a verification code with the type "emailVerification"
func (r *mongoRepository) CheckEmailVerificationCode(ctx context.Context, code string) error {
	err := r.db.C(codeCollection).Remove(bson.M{"type": "emailVerification", "code": code})
	return err
}

// UpdatePhoneVerified takes a user's id and marks their phoneVerified property as true
func (r *mongoRepository) UpdatePhoneVerified(ctx context.Context, id string, value bool) error {
	err := r.db.C(userCollection).Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"phoneVerified": value}})
	return err
}

// UpdateEmailVerified takes a user's id and marks their emailVerified property as true
func (r *mongoRepository) UpdateEmailVerified(ctx context.Context, id string, value bool) error {
	err := r.db.C(userCollection).Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"emailVerified": value}})
	return err
}
