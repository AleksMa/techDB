package repository

import "github.com/AleksMa/techDB/models"

type Repo interface {
	PutForum(forum *models.Forum, ownerID int64) (uint64, error)
	PutThread(thread *models.Thread, forumID int64, authorID int64) (uint64, error)
	PutUser(user *models.User) (uint64, error)
	GetUserByNickname(nickname string) (models.User, error)
	GetUserByID(id int64) (models.User, error)
	GetForumBySlug(slug string) (models.Forum, int64, error)
	GetThreadsByForum(forumID int64) (models.Threads, error)
}
