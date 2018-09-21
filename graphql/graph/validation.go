package graph

import (
	"context"
	"regexp"
	"unicode"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/gqlerror"
	validator "gopkg.in/go-playground/validator.v9"
)

type validation struct {
	validate *validator.Validate
}

// NewValidator initiates a new validator and registers our custom validator tags
func NewValidator() *validation {
	validate := validator.New()
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("complexity", validateComplexity)
	return &validation{validate}
}

// RegisterArgs is a struct representing the arguments that get sent in a Register GraphQL mutation, validation rules can be set as struct tags
type RegisterArgs struct {
	FirstName    string `json:"firstName" validate:"required,min=2,max=32"`
	LastName     string `json:"lastName" validate:"required,min=2,max=32"`
	Email        string `json:"email" validate:"required,email,min=2,max=32"`
	Phone        string `json:"phone" validate:"required,phone,min=2,max=32"`
	Password     string `json:"password" validate:"required,complexity,eqfield=PasswordConf,min=2,max=32"`
	PasswordConf string `json:"passwordConf" validate:"required,complexity,eqfield=Password,min=2,max=32"`
}

// LoginArgs is a struct representing the arguments that get sent in a Login GraphQL mutation, validation rules can be set as struct tags
type LoginArgs struct {
	Email    string `json:"email" validate:"required,email,min=2,max=32""`
	Password string `json:"password" validate:"required,min=2,max=32""`
}

// validatePhone is a handler for a custom validator tag that validates phone numbers
func validatePhone(fl validator.FieldLevel) bool {
	r, _ := regexp.Compile(`(?m)^\s*(?:\+?(\d{1,3}))?[-. (]*(\d{3})[-. )]*(\d{3})[-. ]*(\d{4})(?: *x(\d+))?\s*$`)
	return r.MatchString(fl.Field().String())
}

// validateComplexity is a handler for a custom validator tag that validates password complexity requirements. (1+ uppercase, 1+ numbers, 1+ symbol, & 8 - 32 characters in length)
func validateComplexity(fl validator.FieldLevel) bool {
	fls := fl.Field().String()
	var number, upper, special bool
	// return false if length isn't met
	length := len(fls) >= 8 && len(fls) <= 32
	if length != true {
		return false
	}
	// return false if complexity requirements aren't met
	for _, s := range fls {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
		}
	}
	return number && upper && special
}

// formatValidationErrors is a convienice function that attaches multiple validation errors to a GraphQL response via request context
func formatValidationErrors(ctx context.Context, valErr error) error {
	for _, err := range valErr.(validator.ValidationErrors) {
		graphql.AddError(ctx, gqlerror.Errorf(`Error: Validation for field '`+err.Field()+`' failed.`))
	}
	return nil
}
