package repository

import (
	"fmt"
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

	PutForum(forum *models.Forum) (uint64, *models.Error)
	GetForumBySlug(slug string) (models.Forum, *models.Error)
	GetForumByID(id int64) (models.Forum, *models.Error)

	PutThread(thread *models.Thread) (uint64, *models.Error)
	GetThreadBySlug(slug string) (models.Thread, *models.Error)
	GetThreadByID(id int64) (models.Thread, *models.Error)

	GetThreadsByForum(forumID int64) (models.Threads, error)
	GetStatus() (models.Status, error)
	ReloadDB() error

	PutPost(post *models.Post) (uint64, error)
	UpdateThreadWithID(thread *models.Thread) error
	UpdateThreadWithSlug(thread *models.Thread) error
	GetPostsByThreadID(threadID int64) (models.Posts, error)
	PutVote(vote *models.Vote) (uint64, error)
	GetUsersByForum(forumID int64) (models.Users, error)
	ChangePost(post *models.Post) error
	GetPost(ID int64) (models.Post, error)
}

type DBStore struct {
	DB *pgx.ConnPool
}

func NewDBStore(db *pgx.ConnPool) Repo {
	return &DBStore{
		db,
	}
}

func (store *DBStore) PutThread(thread *models.Thread) (uint64, *models.Error) {
	fmt.Println(thread)
	var ID uint64

	insertQuery := `INSERT INTO threads (created, forumid, message, slug, title, authorid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		thread.Created, thread.ForumID, thread.Message, thread.Slug, thread.Title, thread.AuthorID)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return ID, nil
}

func (store *DBStore) GetThreadsByForum(forumID int64) (models.Threads, error) {
	threads := models.Threads{}

	selectStr := "SELECT ID, created, message, slug, title, authorid FROM threads WHERE forumID = $1"
	rows, err := store.DB.Query(selectStr, forumID)
	if err != nil {
		return threads, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
	}

	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(&thread.ID, &thread.Created, &thread.Message, &thread.Slug, &thread.Title, &thread.AuthorID)
		if err != nil {
			return threads, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
		}
		threads = append(threads, thread)
	}

	rows.Close()

	if err != nil {
		return threads, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

	return threads, nil
}

func (store *DBStore) GetThreadBySlug(slug string) (models.Thread, *models.Error) {
	thread := &models.Thread{}

	selectStr := "SELECT ID, created, forumid, message, slug, title, authorid FROM threads WHERE slug = $1"
	row := store.DB.QueryRow(selectStr, slug)

	err := row.Scan(&thread.ID, &thread.Created, &thread.ForumID, &thread.Message, &thread.Slug, &thread.Title, &thread.AuthorID)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return *thread, models.NewError(http.StatusNotFound, err.Error())
		}
		return *thread, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return *thread, nil
}

func (store *DBStore) GetThreadByID(id int64) (models.Thread, *models.Error) {
	thread := &models.Thread{}

	selectStr := "SELECT ID, created, forumid, message, slug, title, authorid FROM threads WHERE id = $1"
	row := store.DB.QueryRow(selectStr, id)

	err := row.Scan(&thread.ID, &thread.Created, &thread.ForumID, &thread.Message, &thread.Slug, &thread.Title, &thread.AuthorID)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return *thread, models.NewError(http.StatusNotFound, err.Error())
		}
		return *thread, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return *thread, nil
}

func (store *DBStore) UpdateThreadWithSlug(thread *models.Thread) error {
	fmt.Println(thread)

	insertQuery := `UPDATE threads SET message=$1, title=$2 WHERE slug=$3`
	_, err := store.DB.Exec(insertQuery,
		thread.Message, thread.Title, thread.Slug)

	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not put thread: "+err.Error())
	}

	return nil
}

func (store *DBStore) UpdateThreadWithID(thread *models.Thread) error {
	fmt.Println(thread)

	insertQuery := `UPDATE threads SET message=$1, title=$2 WHERE id=$3`
	_, err := store.DB.Exec(insertQuery,
		thread.Message, thread.Title, thread.ID)

	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not put thread: "+err.Error())
	}

	return nil
}

func (store *DBStore) PutPost(post *models.Post) (uint64, error) {
	fmt.Println(post)
	var ID uint64

	insertQuery := `INSERT INTO posts (created, forumid, isedited, message, parentid, authorid, threadid) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		post.Created, post.ForumID, false, post.Message, 0, post.AuthorID, post.Thread)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not put post: "+err.Error())
	}

	return ID, nil
}

func (store *DBStore) GetPostsByThreadID(threadID int64) (models.Posts, error) {
	posts := models.Posts{}

	selectStr := `SELECT ID, created, forumid, isEdited, message, parentid, authorid, threadid 
			FROM posts WHERE threadID = $1`
	rows, err := store.DB.Query(selectStr, threadID)
	if err != nil {
		return posts, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
	}

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Created, &post.ForumID, &post.IsEdited,
			&post.Message, &post.Parent, &post.AuthorID, &post.Thread)
		if err != nil {
			return posts, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
		}
		posts = append(posts, post)
	}

	rows.Close()

	if err != nil {
		return posts, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

	return posts, nil
}

func (store *DBStore) PutVote(vote *models.Vote) (uint64, error) {
	fmt.Println(vote)
	var ID uint64
	positive := vote.Voice == 1

	insertQuery := `INSERT INTO votes (voice, threadid, authorid) VALUES ($1, $2, $3) RETURNING id`
	rows := store.DB.QueryRow(insertQuery, positive, vote.ThreadID, vote.AuthorID)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not put vote: "+err.Error())
	}

	return ID, nil
}

func (store *DBStore) ChangePost(post *models.Post) error {
	fmt.Println(post)

	insertQuery := `UPDATE posts SET message=$1, isedited=$2 WHERE id=$2`
	_, err := store.DB.Exec(insertQuery, post.Message, true, post.ID)

	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not put post: "+err.Error())
	}

	return nil
}

func (store *DBStore) GetPost(id int64) (models.Post, error) {
	post := &models.Post{}

	selectStr := "SELECT ID, created, forumid, isEdited, message, parentid, authorid, threadid FROM posts WHERE id = $1"
	row := store.DB.QueryRow(selectStr, id)

	err := row.Scan(&post.ID, &post.Created, &post.ForumID, &post.IsEdited,
		&post.Message, &post.Parent, &post.AuthorID, &post.Thread)

	if err != nil {
		return *post, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

	return *post, nil
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
