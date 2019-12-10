package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/AleksMa/techDB/repository"
	"time"
)

type UseCase interface {
	PutForum(newForum *models.Forum) error
	PutThread(newThread *models.Thread) error
	PutUser(user *models.User) error
	GetUserByNickname(nickname string) (models.User, error)
	GetForumBySlug(slug string) (models.Forum, error)
	GetThreadsByForum(slug string) (models.Threads, error)
	ChangeUser(user *models.User) error
	GetStatus() (models.Status, error)
	RemoveAllData() error
	PutPost(post *models.Post, threadID int64) (models.Post, error)
	PutPostWithSlug(post *models.Post, threadSlug string) (models.Post, error)
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

func (u *useCase) ChangeUser(user *models.User) error {
	fmt.Println(user)
	//TODO: contains check
	u.repository.ChangeUser(user)
	//TODO: error check
	return nil
}

func (u *useCase) GetStatus() (models.Status, error) {
	return u.repository.GetStatus()
}

func (u *useCase) RemoveAllData() error {
	return u.repository.ReloadDB()
}

func (u *useCase) GetThreadBySlug(slug string) (models.Thread, error) {
	thread, ownerID, _ := u.repository.GetThreadBySlug(slug)
	fmt.Println(thread, ownerID)
	owner, _ := u.repository.GetUserByID(ownerID)
	thread.Author = owner.Nickname
	forum, _, _ := u.repository.GetForumByID(thread.ForumID)
	thread.Forum = forum.Slug
	return thread, nil
}

func (u *useCase) GetThreadByID(id int64) (models.Thread, error) {
	thread, ownerID, _ := u.repository.GetThreadByID(id)
	fmt.Println(thread, ownerID)
	owner, _ := u.repository.GetUserByID(ownerID)
	thread.Author = owner.Nickname
	forum, _, _ := u.repository.GetForumByID(thread.ForumID)
	thread.Forum = forum.Slug
	return thread, nil
}

func (u *useCase) PutPost(post *models.Post, threadID int64) (models.Post, error) {
	fmt.Println(post)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(post.Author)
	// rep return array of ids
	created := time.Now()
	thread, _, _ := u.repository.GetThreadByID(threadID)
	post.Thread = threadID
	post.Forum = thread.Forum
	post.ForumID = thread.ForumID
	post.Created = created
	post.AuthorID = user.ID

	fmt.Println(post)

	u.repository.PutPost(post)
	//TODO: error check
	return *post, nil
}

func (u *useCase) PutPostWithSlug(post *models.Post, threadSlug string) (models.Post, error) {
	fmt.Println(post)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(post.Author)
	// rep return array of ids
	created := time.Now()
	thread, _, _ := u.repository.GetThreadBySlug(threadSlug)
	post.Thread = thread.ID
	post.Forum = thread.Forum
	post.AuthorID = user.ID
	post.ForumID = thread.ForumID

	post.Created = created

	fmt.Println(post)

	u.repository.PutPost(post)
	//TODO: error check
	return *post, nil
}
