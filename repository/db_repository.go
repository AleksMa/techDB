package repository

import (
	"github.com/AleksMa/techDB/models"
	"github.com/jackc/pgx"
)

type Repo interface {
	PutUser(user *models.User) (uint64, *models.Error)
	GetDupUsers(user *models.User) (models.Users, *models.Error)
	GetUserByID(id int64) (models.User, *models.Error)
	GetUserByNickname(nickname string) (models.User, *models.Error)
	ChangeUser(user *models.User) *models.Error

	PutForum(forum *models.Forum) (uint64, *models.Error)
	GetForumBySlug(slug string) (models.Forum, *models.Error)
	GetForumByID(id int64) (models.Forum, *models.Error)

	PutThread(thread *models.Thread) (uint64, *models.Error)
	GetThreadBySlug(slug string) (models.Thread, *models.Error)
	GetThreadByID(id int64) (models.Thread, *models.Error)
	GetThreadsByForum(forumID int64, params models.ThreadParams) (models.Threads, *models.Error)

	PutPost(post *models.Post) (uint64, *models.Error)
	GetPost(ID int64) (models.Post, *models.Error)

	UpdateThreadWithID(thread *models.Thread) error
	UpdateThreadWithSlug(thread *models.Thread) error
	GetPostsByThreadID(threadID int64) (models.Posts, error)
	PutVote(vote *models.Vote) (uint64, error)
	GetUsersByForum(forumID int64) (models.Users, error)
	ChangePost(post *models.Post) error

	GetStatus() (models.Status, error)
	ReloadDB() error
}

type DBStore struct {
	DB *pgx.ConnPool
}

func NewDBStore(db *pgx.ConnPool) Repo {
	return &DBStore{
		db,
	}
}

func (store *DBStore) GetStatus() (models.Status, error) {
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

func (store *DBStore) ReloadDB() error {
	_, err := store.DB.Exec(models.InitScript)
	return err
}
