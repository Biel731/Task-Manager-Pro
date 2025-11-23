package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(input RegisterInput) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	if err := CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func Login(input LoginInput) (*User, error) {
	user, err := FindUserByEmail(input.Email)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("Invalid email or password")
	}

	return user, nil
}
