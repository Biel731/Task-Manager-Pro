package users

import (
	"github.com/bielrodrigues/task-manager-pro-backend/internal/database"
)

func CreateUser(user *User) error {
	return database.DB.Create(user).Error
}

func FindUserByEmail(email string) (*User, error) {
	var user User
	err := database.DB.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
