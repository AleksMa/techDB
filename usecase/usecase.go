package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/AleksMa/techDB/repository"
	"time"
)

type UseCase interface {
	PutUser(user *models.User) (models.Users, *models.Error)
	GetUserByNickname(nickname string) (models.User, *models.Error)
	ChangeUser(user *models.User) (models.User, *models.Error)

	PutForum(newForum *models.Forum) error
	PutThread(newThread *models.Thread) error
	GetForumBySlug(slug string) (models.Forum, error)
	GetThreadsByForum(slug string) (models.Threads, error)
	GetStatus() (models.Status, error)
	RemoveAllData() error
	PutPost(post *models.Post, threadID int64) (models.Post, error)
	PutPostWithSlug(post *models.Post, threadSlug string) (models.Post, error)
	GetThreadByID(id int64) (models.Thread, error)
	GetThreadBySlug(slug string) (models.Thread, error)
	UpdateThreadWithID(thread *models.Thread) (models.Thread, error)
	UpdateThreadWithSlug(thread *models.Thread) (models.Thread, error)
	GetPostsByThreadID(id int64) (models.Posts, error)
	GetPostsByThreadSlug(slug string) (models.Posts, error)
	PutVote(vote *models.Vote) (models.Vote, error)
	PutVoteWithSlug(vote *models.Vote, slug string) (models.Vote, error)
	GetUsersByForum(slug string) (models.Users, error)
	ChangePost(post *models.Post) error
	GetPostFull(id int64) (models.PostFull, error)
}

type useCase struct {
	repository repository.Repo
}

func NewUseCase(repo repository.Repo) UseCase {
	return &useCase{
		repository: repo,
	}
}

func (u *useCase) PutUser(user *models.User) (models.Users, *models.Error) {
	fmt.Println(user)

	if err := user.Validate(); err != nil {
		return nil, err
	}

	users, _ := u.repository.GetDupUsers(user)
	if users != nil && len(users) != 0 {
		fmt.Println("DUP: ", users)
		return users, nil
	}

	_, err := u.repository.PutUser(user)
	return nil, err
}

func (u *useCase) GetUserByNickname(nickname string) (models.User, *models.Error) {
	return u.repository.GetUserByNickname(nickname)
}

func (u *useCase) ChangeUser(user *models.User) (models.User, *models.Error) {
	tempUser, err := u.repository.GetUserByNickname(user.Nickname)
	if err != nil {
		return *user, err
	}
	err = u.repository.ChangeUser(user)
	fmt.Println(*user)
	fmt.Println(tempUser)
	return *user, err
}

func (u *useCase) GetUsersByForum(slug string) (models.Users, error) {
	forum, _, _ := u.repository.GetForumBySlug(slug)

	users, _ := u.repository.GetUsersByForum(forum.ID)
	fmt.Println(users)
	for i, _ := range users {
		user, _ := u.repository.GetUserByID(users[i].ID)
		users[i].Nickname = user.Nickname
		users[i].Fullname = user.Fullname
		users[i].About = user.About
		users[i].Email = user.Email
	}
	return users, nil
}

func (u *useCase) PutForum(newForum *models.Forum) error {
	fmt.Println(newForum)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(newForum.Owner)
	u.repository.PutForum(newForum, user.ID)
	//TODO: error check
	return nil
}

func (u *useCase) GetForumBySlug(slug string) (models.Forum, error) {
	forum, ownerID, _ := u.repository.GetForumBySlug(slug)
	fmt.Println(forum, ownerID)
	owner, _ := u.repository.GetUserByID(ownerID)
	forum.Owner = owner.Nickname
	return forum, nil
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

func (u *useCase) UpdateThreadWithID(thread *models.Thread) (models.Thread, error) {
	fmt.Println(thread)
	//TODO: contains check
	u.repository.UpdateThreadWithID(thread)
	//TODO: error check
	return *thread, nil
}

func (u *useCase) UpdateThreadWithSlug(thread *models.Thread) (models.Thread, error) {
	fmt.Println(thread)
	//TODO: contains check
	u.repository.UpdateThreadWithSlug(thread)
	//TODO: error check
	return *thread, nil
}

func (u *useCase) GetPostsByThreadID(id int64) (models.Posts, error) {
	thread, _, _ := u.repository.GetThreadByID(id)

	posts, _ := u.repository.GetPostsByThreadID(thread.ID)
	for i, _ := range posts {
		posts[i].Thread = thread.ID
		user, _ := u.repository.GetUserByID(posts[i].AuthorID)
		posts[i].Author = user.Nickname
	}
	return posts, nil
}

func (u *useCase) GetPostsByThreadSlug(slug string) (models.Posts, error) {
	thread, _, _ := u.repository.GetThreadBySlug(slug)

	posts, _ := u.repository.GetPostsByThreadID(thread.ID)
	for i, _ := range posts {
		posts[i].Thread = thread.ID
		user, _ := u.repository.GetUserByID(posts[i].AuthorID)
		posts[i].Author = user.Nickname
	}
	return posts, nil
}

func (u *useCase) PutVote(vote *models.Vote) (models.Vote, error) {
	fmt.Println(vote)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(vote.Nickname)
	vote.AuthorID = user.ID

	fmt.Println(vote)

	u.repository.PutVote(vote)
	//TODO: error check
	return *vote, nil
}

func (u *useCase) ChangePost(post *models.Post) error {
	fmt.Println(post)
	//TODO: contains check
	u.repository.ChangePost(post)
	//TODO: error check
	return nil
}

func (u *useCase) GetPostFull(id int64) (models.PostFull, error) {
	var postFull models.PostFull
	var err error
	postFull.Post, err = u.repository.GetPost(id)
	fmt.Println(postFull)
	fmt.Println(err)

	postFull.Author, _ = u.repository.GetUserByID(postFull.Post.AuthorID)
	postFull.Thread, _, _ = u.repository.GetThreadByID(postFull.Post.Thread)
	postFull.Forum, _, _ = u.repository.GetForumByID(postFull.Post.ForumID)

	return postFull, nil
}

func (u *useCase) PutVoteWithSlug(vote *models.Vote, slug string) (models.Vote, error) {
	fmt.Println(vote)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(vote.Nickname)
	vote.AuthorID = user.ID
	thread, _, _ := u.repository.GetThreadBySlug(slug)
	vote.ThreadID = thread.ID
	fmt.Println(vote)

	u.repository.PutVote(vote)
	//TODO: error check
	return *vote, nil
}

func (u *useCase) GetStatus() (models.Status, error) {
	return u.repository.GetStatus()
}

func (u *useCase) RemoveAllData() error {
	return u.repository.ReloadDB()
}
