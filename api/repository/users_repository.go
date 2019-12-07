package repository

import "github.com/CoolCodeTeam/2019_2_CoolCodeMicroServices/utils/models"

//go:generate moq -out user_repo_mock.go . UserRepo

type UserRepo interface {
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(ID uint64) (models.User, error)
	PutUser(newUser *models.User) (uint64, error)
	Replace(ID uint64, newUser *models.User) error
	Contains(user models.User) bool
	GetUsers() (models.Users, error)
}
