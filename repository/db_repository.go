package repository

import (
	"database/sql"
	"fmt"
	"github.com/AleksMa/techDB/models"
	"net/http"
)

type DBStore struct {
	DB *sql.DB
}

func NewDBStore(db *sql.DB) Repo {
	return &DBStore{
		db,
	}
}

func (store *DBStore) PutForum(forum *models.Forum, ownerID int64) (uint64, error) {
	fmt.Println(forum)
	var ID uint64

	insertQuery := `INSERT INTO forums (slug, title, authorid) VALUES ($1, $2, $3) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		forum.Slug, forum.Title, ownerID)

	err := rows.Scan(&ID)
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not put user: "+err.Error())
	}

	return ID, nil
}

func (store *DBStore) PutThread(thread *models.Thread, forumID int64, authorID int64) (uint64, error) {
	fmt.Println(thread)
	var ID uint64

	insertQuery := `INSERT INTO threads (created, forumid, message, slug, title, authorid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		thread.Created, forumID, thread.Message, thread.Slug, thread.Title, authorID)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not put thread: "+err.Error())
	}

	return ID, nil
}

func (store *DBStore) PutUser(user *models.User) (uint64, error) {
	fmt.Println(user)
	var ID uint64

	insertQuery := `INSERT INTO users (nickname, about, email, fullname) VALUES ($1, $2, $3, $4) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		user.Nickname, user.About, user.Email, user.Fullname)

	err := rows.Scan(&ID)
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not put user: "+err.Error())
	}

	return ID, nil
}

func (store *DBStore) GetUserByNickname(nickname string) (models.User, error) {
	user := &models.User{}

	selectStr := "SELECT id, nickname, about, email, fullname FROM users WHERE nickname = $1"
	row := store.DB.QueryRow(selectStr, nickname)

	err := row.Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Fullname)

	if err != nil {
		return *user, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

	return *user, nil
}

func (store *DBStore) GetUserByID(id int64) (models.User, error) {
	user := &models.User{}

	selectStr := "SELECT id, nickname, about, email, fullname FROM users WHERE id = $1"
	row := store.DB.QueryRow(selectStr, id)

	err := row.Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Fullname)

	if err != nil {
		return *user, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

	return *user, nil
}

func (store *DBStore) GetForumBySlug(slug string) (models.Forum, int64, error) {
	forum := &models.Forum{}
	var ownerID int64

	selectStr := "SELECT ID, slug, title, authorid FROM forums WHERE slug = $1"
	row := store.DB.QueryRow(selectStr, slug)

	err := row.Scan(&forum.ID, &forum.Slug, &forum.Title, &ownerID)

	if err != nil {
		return *forum, 0, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

	return *forum, ownerID, nil
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
