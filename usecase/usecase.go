package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/AleksMa/techDB/repository"
	"net/http"
	"time"
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
	GetUserByID(id int64) (models.User, *models.Error)
	GetThreadsByForum(slug string, params models.ThreadParams) (models.Threads, *models.Error)

	GetStatus() (models.Status, error)
	RemoveAllData() error
	PutPost(post *models.Post, threadID int64) (models.Post, error)
	PutPostWithSlug(post *models.Post, threadSlug string) (models.Post, error)
	GetThreadByID(id int64) (models.Thread, error)
	UpdateThreadWithID(thread *models.Thread) (models.Thread, error)
	UpdateThreadWithSlug(thread *models.Thread) (models.Thread, error)
	GetPostsByThreadID(id int64) (models.Posts, error)
	GetPostsByThreadSlug(slug string) (models.Posts, error)
	PutVote(vote *models.Vote) (models.Vote, error)
	PutVoteWithSlug(vote *models.Vote, slug string) (models.Vote, error)
	ChangePost(post *models.Post) error
	GetPostFull(id int64) (models.PostFull, error)

	GetUsersByForum(slug string) (models.Users, error)
}

type useCase struct {
	repository repository.Repo
}

func NewUseCase(repo repository.Repo) UseCase {
	return &useCase{
		repository: repo,
	}
}

func (u *useCase) PutThread(thread *models.Thread) (models.Thread, *models.Error) {
	fmt.Println(thread)
	if thread.Slug != "" {
		dupForum, err := u.GetThreadBySlug(thread.Slug)
		if err == nil || err.Code != http.StatusNotFound {
			fmt.Println("DUP: ", dupForum)
			return dupForum, models.NewError(http.StatusConflict, "forum already created")
		}
	}

	user, err := u.GetUserByNickname(thread.Author)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return *thread, models.NewError(http.StatusNotFound, "No user found: "+err.Message)
		}
		return *thread, models.NewError(http.StatusInternalServerError, err.Message)
	}

	forum, err := u.GetForumBySlug(thread.Forum)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return *thread, models.NewError(http.StatusNotFound, "No forum found: "+err.Message)
		}
		return *thread, models.NewError(http.StatusInternalServerError, err.Message)
	}
	thread.AuthorID = user.ID
	thread.ForumID = forum.ID

	id, err := u.repository.PutThread(thread)
	if err != nil {
		return *thread, err
	}
	thread.ID = int64(id)
	return *thread, nil
}

func (u *useCase) GetThreadBySlug(slug string) (models.Thread, *models.Error) {
	thread, err := u.repository.GetThreadBySlug(slug)
	fmt.Println(thread)

	owner, err := u.GetUserByID(thread.AuthorID)
	if err != nil {
		return thread, err
	}
	thread.Author = owner.Nickname

	forum, err := u.repository.GetForumBySlug(slug)
	if err != nil {
		return thread, err
	}

	thread.Forum = forum.Slug
	return thread, nil
}

func (u *useCase) GetThreadsByForum(slug string, params models.ThreadParams) (models.Threads, *models.Error) {
	forum, err := u.GetForumBySlug(slug)
	if err != nil {
		return nil, err
	}

	threads, _ := u.repository.GetThreadsByForum(forum.ID, params)
	for i, _ := range threads {
		threads[i].Forum = forum.Slug
		user, _ := u.repository.GetUserByID(threads[i].AuthorID)
		threads[i].Author = user.Nickname
	}
	return threads, nil
}

func (u *useCase) GetThreadByID(id int64) (models.Thread, error) {
	thread, _ := u.repository.GetThreadByID(id)
	fmt.Println(thread)
	owner, _ := u.repository.GetUserByID(thread.AuthorID)
	thread.Author = owner.Nickname
	forum, _ := u.GetForumByID(thread.ForumID)
	thread.Forum = forum.Slug
	return thread, nil
}

func (u *useCase) PutPost(post *models.Post, threadID int64) (models.Post, error) {
	fmt.Println(post)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(post.Author)
	// rep return array of ids
	created := time.Now()
	thread, _ := u.repository.GetThreadByID(threadID)
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
	thread, _ := u.repository.GetThreadBySlug(threadSlug)
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
	thread, _ := u.repository.GetThreadByID(id)

	posts, _ := u.repository.GetPostsByThreadID(thread.ID)
	for i, _ := range posts {
		posts[i].Thread = thread.ID
		user, _ := u.repository.GetUserByID(posts[i].AuthorID)
		posts[i].Author = user.Nickname
	}
	return posts, nil
}

func (u *useCase) GetPostsByThreadSlug(slug string) (models.Posts, error) {
	thread, _ := u.repository.GetThreadBySlug(slug)

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
	postFull.Thread, _ = u.repository.GetThreadByID(postFull.Post.Thread)
	postFull.Forum, _ = u.GetForumByID(postFull.Post.ForumID)

	return postFull, nil
}

func (u *useCase) PutVoteWithSlug(vote *models.Vote, slug string) (models.Vote, error) {
	fmt.Println(vote)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(vote.Nickname)
	vote.AuthorID = user.ID
	thread, _ := u.repository.GetThreadBySlug(slug)
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
