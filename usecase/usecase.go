package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/AleksMa/techDB/repository"
)

type UseCase interface {
	PutForum(newForum *models.Forum) error
	PutThread(newThread *models.Thread) error
	PutUser(user *models.User) error
	GetUserByNickname(nickname string) (models.User, error)
	GetForumBySlug(slug string) (models.Forum, error)
	GetThreadsByForum(slug string) (models.Threads, error)
}

type useCase struct {
	repository repository.Repo
}

func NewUseCase(repo repository.Repo) UseCase {
	return &useCase{
		repository: repo,
	}
}

func (u *useCase) PutForum(newForum *models.Forum) error {
	fmt.Println(newForum)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(newForum.Owner)
	u.repository.PutForum(newForum, user.ID)
	//TODO: error check
	return nil
}

func (u *useCase) PutThread(newThread *models.Thread) error {
	fmt.Println(newThread)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(newThread.Author)
	forum, _, _ := u.repository.GetForumBySlug(newThread.Forum)

	fmt.Println(user, forum)

	u.repository.PutThread(newThread, forum.ID, user.ID)
	//TODO: error check
	return nil
}

func (u *useCase) PutUser(user *models.User) error {
	fmt.Println(user)
	//TODO: contains check
	u.repository.PutUser(user)
	//TODO: error check
	return nil
}

func (u *useCase) GetUserByNickname(nickname string) (models.User, error) {
	return u.repository.GetUserByNickname(nickname)
}

func (u *useCase) GetForumBySlug(slug string) (models.Forum, error) {
	forum, ownerID, _ := u.repository.GetForumBySlug(slug)
	fmt.Println(forum, ownerID)
	owner, _ := u.repository.GetUserByID(ownerID)
	forum.Owner = owner.Nickname
	return forum, nil
}

func (u *useCase) GetThreadsByForum(slug string) (models.Threads, error) {
	forum, _, _ := u.repository.GetForumBySlug(slug)

	threads, _ := u.repository.GetThreadsByForum(forum.ID)
	for i, _ := range threads {
		threads[i].Forum = forum.Slug
		user, _ := u.repository.GetUserByID(threads[i].AuthorID)
		threads[i].Author = user.Nickname
	}
	return threads, nil
}
