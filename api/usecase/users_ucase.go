package useCase

import (
	"../models"
	"../repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type UsersUseCase interface {
	GetUserByID(id uint64) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	SignUp(user *models.User) error
	Login(user models.User) (models.User, error)
	ChangeUser(user *models.User) error
	FindUsers(name string) (models.Users, error)
}

type usersUseCase struct {
	repository repository.UserRepo
}

func NewUserUseCase(repo repository.UserRepo) UsersUseCase {
	return &usersUseCase{
		repository: repo,
	}
}

func (u *usersUseCase) GetUserByID(id uint64) (models.User, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		return user, err
	}
	if !u.Valid(user) {
		return user, fmt.Errorf("no such user")//models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
	return user, nil
}

func (u *usersUseCase) GetUserByEmail(email string) (models.User, error) {
	return u.repository.GetUserByEmail(email)
}

func (u *usersUseCase) SignUp(newUser *models.User) error {
	if u.repository.Contains(*newUser) {
		return fmt.Errorf("already exist")
	} else {
		if newUser.Name == "" {
			newUser.Name = "John Doe"
		}
		newUser.Password = hashAndSalt(newUser.Password)
		_, err := u.repository.PutUser(newUser)
		if err != nil { // return 500 Internal Server Error.
			return models.NewServerError(err, http.StatusInternalServerError, "")
		}
	}
	return nil
}
