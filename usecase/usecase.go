package useCase

import (
	"github.com/AleksMa/techDB/models"
	"github.com/AleksMa/techDB/repository"
)

type UseCase interface {
	PutUser(user *models.User) (models.Users, *models.Error)
	GetUserByNickname(nickname string) (models.User, *models.Error)
	ChangeUser(userUpd *models.UpdateUserFields, nickname string) (models.User, *models.Error)

	PutForum(forum *models.Forum) (models.Forum, *models.Error)
	GetForumBySlug(slug string) (models.Forum, *models.Error)
	GetForumByID(id int64) (models.Forum, *models.Error)

	PutThread(thread *models.Thread) (models.Thread, *models.Error)
	GetThreadBySlug(slug string) (models.Thread, *models.Error)
	GetThreadByID(id int64) (models.Thread, *models.Error)
	GetUserByID(id int64) (models.User, *models.Error)
	GetThreadsByForum(slug string, params models.ThreadParams) (models.Threads, *models.Error)
	UpdateThreadWithID(thread *models.Thread) (models.Thread, *models.Error)
	UpdateThreadWithSlug(thread *models.Thread) (models.Thread, *models.Error)

	PutPost(post *models.Post) (*models.Post, *models.Error)
	PutPostWithSlug(post *models.Post, threadSlug string) (*models.Post, *models.Error)
	GetPostFull(id int64, fields []string) (models.PostFull, *models.Error)
	GetPostByID(id int64) (models.Post, *models.Error)
	ChangePost(post *models.Post) (models.Post, *models.Error)

	PutVote(vote *models.Vote) (models.Thread, *models.Error)
	PutVoteWithSlug(vote *models.Vote, slug string) (models.Thread, *models.Error)

	GetPostsByThreadID(id int64) (models.Posts, error)
	GetPostsByThreadSlug(slug string) (models.Posts, error)

	GetUsersByForum(slug string) (models.Users, error)

	GetStatus() (models.Status, error)
	RemoveAllData() error
}

type useCase struct {
	repository repository.Repo
}

func NewUseCase(repo repository.Repo) UseCase {
	return &useCase{
		repository: repo,
	}
}

func (u *useCase) GetStatus() (models.Status, error) {
	return u.repository.GetStatus()
}

func (u *useCase) RemoveAllData() error {
	return u.repository.ReloadDB()
}
