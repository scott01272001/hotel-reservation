package types

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = bcrypt.DefaultCost
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" validate:"omitempty,email"`
}

type UpdateUserParamType UpdateUserParams

type CreateUserParam struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=4"`
}

type User struct {
	ID                string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lastName" json:"lastName"`
	Email             string `bson:"email" json:"email"`
	EncryptedPasswrod string `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params *CreateUserParam) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPasswrod: string(encpw),
	}, nil
}
