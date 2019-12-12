package repository

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/jackc/pgx"
	"net/http"
)

func (store *DBStore) PutForum(forum *models.Forum) (uint64, *models.Error) {
	fmt.Println(forum)
	var ID uint64

	insertQuery := `INSERT INTO forums (slug, title, authorid) VALUES ($1, $2, $3) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		forum.Slug, forum.Title, forum.OwnerID)

	err := rows.Scan(&ID)
	if err != nil {
		return 0, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return ID, nil
}

func (store *DBStore) GetForumBySlug(slug string) (models.Forum, *models.Error) {
	tx, _ := store.DB.Begin()
	defer tx.Rollback()

	forum := &models.Forum{}

	selectStr := "SELECT ID, slug, title, authorid FROM forums WHERE slug = $1"
	row := tx.QueryRow(selectStr, slug)

	err := row.Scan(&forum.ID, &forum.Slug, &forum.Title, &forum.OwnerID)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return *forum, models.NewError(http.StatusNotFound, err.Error())
		}
		return *forum, models.NewError(http.StatusInternalServerError, err.Error())
	}

	selectStr = "SELECT COUNT(*) FROM threads WHERE forumid = $1"
	row = tx.QueryRow(selectStr, forum.ID)
	row.Scan(&forum.Threads)

	selectStr = "SELECT COUNT(*) FROM posts WHERE forumid = $1"
	row = tx.QueryRow(selectStr, forum.ID)
	row.Scan(&forum.Posts)

	tx.Commit()

	return *forum, nil
}

func (store *DBStore) GetForumByID(id int64) (models.Forum, *models.Error) {
	tx, _ := store.DB.Begin()
	defer tx.Rollback()

	forum := &models.Forum{}

	selectStr := "SELECT ID, slug, title, authorid FROM forums WHERE id = $1"
	row := store.DB.QueryRow(selectStr, id)

	err := row.Scan(&forum.ID, &forum.Slug, &forum.Title, &forum.OwnerID)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return *forum, models.NewError(http.StatusNotFound, err.Error())
		}
		return *forum, models.NewError(http.StatusInternalServerError, err.Error())
	}

	selectStr = "SELECT COUNT(*) FROM threads WHERE forumid = $1"
	row = tx.QueryRow(selectStr, forum.ID)
	row.Scan(&forum.Threads)

	selectStr = "SELECT COUNT(*) FROM posts WHERE forumid = $1"
	row = tx.QueryRow(selectStr, forum.ID)
	row.Scan(&forum.Posts)

	tx.Commit()

	return *forum, nil
}
