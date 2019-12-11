package repository

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/jackc/pgx"
	"net/http"
)

func (store *DBStore) PutUser(user *models.User) (uint64, *models.Error) {
	fmt.Println(user)
	var ID uint64

	insertQuery := `INSERT INTO users (nickname, about, email, fullname) VALUES ($1, $2, $3, $4) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		user.Nickname, user.About, user.Email, user.Fullname)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println(err)
		return 0, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return ID, nil
}

func (store *DBStore) GetDupUsers(user *models.User) (models.Users, *models.Error) {
	users := models.Users{}

	selectStr := "SELECT DISTINCT nickname, about, email, fullname FROM users WHERE nickname=$1 OR email=$2"

	rows, err := store.DB.Query(selectStr, user.Nickname, user.Email)
	if err != nil {
		fmt.Println(err)
		return users, models.NewError(http.StatusInternalServerError, err.Error())
	}

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Nickname, &user.About, &user.Email, &user.Fullname)
		if err != nil {
			return users, models.NewError(http.StatusInternalServerError, err.Error())
		}
		users = append(users, user)
	}

	rows.Close()

	if err != nil {
		return users, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return users, nil
}

func (store *DBStore) GetUserByNickname(nickname string) (models.User, *models.Error) {
	user := &models.User{}

	selectStr := "SELECT id, nickname, about, email, fullname FROM users WHERE nickname = $1"
	row := store.DB.QueryRow(selectStr, nickname)

	err := row.Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Fullname)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return *user, models.NewError(http.StatusNotFound, err.Error())
		}
		return *user, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return *user, nil
}

func (store *DBStore) GetUserByID(id int64) (models.User, *models.Error) {
	user := &models.User{}

	selectStr := "SELECT id, nickname, about, email, fullname FROM users WHERE id = $1"
	row := store.DB.QueryRow(selectStr, id)

	err := row.Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Fullname)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return *user, models.NewError(http.StatusNotFound, err.Error())
		}
		return *user, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return *user, nil
}

func (store *DBStore) ChangeUser(user *models.User) *models.Error {

	insertQuery := `UPDATE users SET about=$1, email=$2, fullname=$3 WHERE nickname=$4`
	_, err := store.DB.Exec(insertQuery,
		user.About, user.Email, user.Fullname, user.Nickname)

	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok && pgerr.Code == "23505" {
			return models.NewError(http.StatusConflict, err.Error())
		}
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (store *DBStore) GetUsersByForum(forumID int64) (models.Users, error) {
	users := models.Users{}

	selectStr := "SELECT DISTINCT authorid FROM posts WHERE forumID = $1 UNION SELECT authorid FROM threads WHERE forumID = $1"
	rows, err := store.DB.Query(selectStr, forumID)
	if err != nil {
		fmt.Println(err)
		return users, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
	}

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID)
		if err != nil {
			return users, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
		}
		users = append(users, user)
	}

	rows.Close()

	if err != nil {
		return users, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

	return users, nil
}
