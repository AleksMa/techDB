package repository

import (
	"github.com/AleksMa/techDB/models"
	"github.com/jackc/pgx"
	"net/http"
)

type Repo interface {
	PutUser(user *models.User) (uint64, *models.Error)
	GetDupUsers(user *models.User) (models.Users, *models.Error)
	GetUserByID(id int64) (models.User, *models.Error)
	GetUserByNickname(nickname string) (models.User, *models.Error)
	ChangeUser(user *models.User) *models.Error
	GetUsersByForum(forumID int64, params models.UserParams) (models.Users, *models.Error)

	PutForum(forum *models.Forum) (uint64, *models.Error)
	GetForumBySlug(slug string) (models.Forum, *models.Error)
	GetForumByID(id int64) (models.Forum, *models.Error)

	PutThread(thread *models.Thread) (uint64, *models.Error)
	GetThreadBySlug(slug string) (models.Thread, *models.Error)
	GetThreadByID(id int64) (models.Thread, *models.Error)
	GetThreadsByForum(forumID int64, params models.ThreadParams) (models.Threads, *models.Error)
	UpdateThreadWithID(thread *models.Thread) *models.Error
	UpdateThreadWithSlug(thread *models.Thread) *models.Error

	PutPost(post *models.Post) (uint64, *models.Error)
	GetPost(ID int64) (models.Post, *models.Error)
	ChangePost(post *models.Post) *models.Error
	GetPostsByThreadID(threadID int64, params models.PostParams) (models.Posts, *models.Error)

	UpdateVote(vote *models.Vote) (int, *models.Error)
	PutVote(vote *models.Vote) (uint64, *models.Error)

	GetStatus() (models.Status, *models.Error)
	ReloadDB() *models.Error
}

type DBStore struct {
	DB *pgx.ConnPool
}

func NewDBStore(db *pgx.ConnPool) Repo {
	return &DBStore{
		db,
	}
}

func (store *DBStore) GetStatus() (models.Status, *models.Error) {
	tx, _ := store.DB.Begin()
	defer tx.Rollback()

	status := &models.Status{}

	row := tx.QueryRow(`SELECT count(*) FROM forums`)
	row.Scan(&status.Forum)

	row = tx.QueryRow(`SELECT count(*) FROM posts`)
	row.Scan(&status.Post)

	row = tx.QueryRow(`SELECT count(*) FROM threads`)
	row.Scan(&status.Thread)

	row = tx.QueryRow(`SELECT count(*) FROM users`)
	row.Scan(&status.User)

	tx.Commit()

	return *status, nil
}

func (store *DBStore) ReloadDB() *models.Error {
	_, err := store.DB.Exec(models.InitScript)
	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
