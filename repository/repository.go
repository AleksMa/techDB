package repository

import "github.com/AleksMa/techDB/models"

type Repo interface {
	PutForum(forum *models.Forum, ownerID int64) (uint64, error)
	PutThread(thread *models.Thread, forumID int64, authorID int64) (uint64, error)
	PutUser(user *models.User) (uint64, error)
	GetUserByNickname(nickname string) (models.User, error)
	GetUserByID(id int64) (models.User, error)
	GetForumBySlug(slug string) (models.Forum, int64, error)
	GetForumByID(id int64) (models.Forum, int64, error)
	GetThreadsByForum(forumID int64) (models.Threads, error)
	ChangeUser(user *models.User) error
	GetStatus() (models.Status, error)
	ReloadDB() error
	GetThreadBySlug(slug string) (models.Thread, int64, error)
	GetThreadByID(id int64) (models.Thread, int64, error)
	PutPost(post *models.Post) (uint64, error)
	UpdateThreadWithID(thread *models.Thread) error
	UpdateThreadWithSlug(thread *models.Thread) error
	GetPostsByThreadID(threadID int64) (models.Posts, error)
	//GetPostsByThreadSlug(slug int64) (models.Posts, error)
	PutVote(vote *models.Vote) (uint64, error)
	GetUsersByForum(forumID int64) (models.Users, error)
	ChangePost(post *models.Post) error
	GetPost(ID int64) (models.Post, error)
}
