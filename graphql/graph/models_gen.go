// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	fmt "fmt"
	io "io"
	strconv "strconv"
)

type LoginArgs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterArgs struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
}
type Session struct {
	ID    string `json:"id"`
	Roles []Role `json:"roles"`
}
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleOwner Role = "OWNER"
	RoleUser  Role = "USER"
)

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleOwner, RoleUser:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
