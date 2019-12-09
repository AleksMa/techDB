package useCase

import (
	"github.com/AleksMa/techDB/repository"
)

type UseCase interface {
}

type useCase struct {
	repository repository.Repo
}

func NewUseCase(repo repository.Repo) UseCase {
	return &useCase{
		repository: repo,
	}
}
