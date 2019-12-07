package delivery

import (
	"../usecase"
)

type UserHandlers struct {
	Users    useCase.UsersUseCase
	//utils    utils.HandlersUtils
}

func NewUsersHandlers(users useCase.UsersUseCase) *UserHandlers {
	return &UserHandlers{
		Users:    users,
		//utils:    utils,
	}
}
