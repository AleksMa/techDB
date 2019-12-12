package repository

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/jackc/pgx"
	"net/http"
	"strconv"
)

func (store *DBStore) PutThread(thread *models.Thread) (uint64, *models.Error) {
	fmt.Println(thread)
	var ID uint64

	var insertQuery string
	params := []interface{}{thread.Created, thread.ForumID, thread.Message, thread.Title, thread.AuthorID}
	if thread.Slug == "" {
		insertQuery = `INSERT INTO threads (created, forumid, message, title, authorid) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	} else {
		insertQuery = `INSERT INTO threads (created, forumid, message, title, authorid, slug) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		params = append(params, thread.Slug)
	}

	//insertQuery := `INSERT INTO threads (created, forumid, message, slug, title, authorid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	rows := store.DB.QueryRow(insertQuery, params...)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return ID, nil
}

func (store *DBStore) GetThreadsByForum(forumID int64, params models.ThreadParams) (models.Threads, *models.Error) {
	fmt.Println("PARAMS: ", params)
	threads := models.Threads{}
	args := 1
	curParams := []interface{}{forumID}

	selectStr := "SELECT ID, created, message, slug, title, authorid, vote FROM threads WHERE forumID = $1"

	if !params.Since.IsZero() {
		selectStr += " AND created"
		if params.Desc {
			selectStr += " <="
		} else {
			selectStr += " >="
		}
		selectStr += " $2"
		args++
		curParams = append(curParams, params.Since)
	}
	selectStr += " ORDER BY created"
	if params.Desc {
		selectStr += " DESC"
	}
	if params.Limit != -1 {
		selectStr += " LIMIT $"
		selectStr += strconv.Itoa(args + 1)
		args++
		curParams = append(curParams, params.Limit)
	}
	selectStr += ";"

	fmt.Println(selectStr)

	rows, err := store.DB.Query(selectStr, curParams...)
	if err != nil {
		return threads, models.NewError(http.StatusInternalServerError, err.Error())
	}

	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(&thread.ID, &thread.Created, &thread.Message, &thread.Slug, &thread.Title, &thread.AuthorID, &thread.Votes)
		if err != nil {
			return threads, models.NewError(http.StatusInternalServerError, err.Error())
		}
		threads = append(threads, thread)
	}

	rows.Close()

	if err != nil {
		return threads, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return threads, nil
}

func (store *DBStore) GetThreadBySlug(slug string) (models.Thread, *models.Error) {
	thread := &models.Thread{}

	selectStr := "SELECT ID, created, forumid, message, slug, title, authorid, vote FROM threads WHERE slug = $1"
	row := store.DB.QueryRow(selectStr, slug)

	err := row.Scan(&thread.ID, &thread.Created, &thread.ForumID, &thread.Message, &thread.Slug, &thread.Title, &thread.AuthorID, &thread.Votes)

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

	selectStr := "SELECT ID, created, forumid, message, slug, title, authorid, vote FROM threads WHERE id = $1"
	row := store.DB.QueryRow(selectStr, id)

	err := row.Scan(&thread.ID, &thread.Created, &thread.ForumID, &thread.Message, &thread.Slug, &thread.Title, &thread.AuthorID, &thread.Votes)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return *thread, models.NewError(http.StatusNotFound, err.Error())
		}
		return *thread, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return *thread, nil
}

func (store *DBStore) UpdateThreadWithSlug(thread *models.Thread) *models.Error {
	fmt.Println(thread)

	insertQuery := `UPDATE threads SET message=$1, title=$2 WHERE slug=$3`
	_, err := store.DB.Exec(insertQuery,
		thread.Message, thread.Title, thread.Slug)

	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (store *DBStore) UpdateThreadWithID(thread *models.Thread) *models.Error {
	fmt.Println(thread)

	insertQuery := `UPDATE threads SET message=$1, title=$2 WHERE id=$3`
	_, err := store.DB.Exec(insertQuery,
		thread.Message, thread.Title, thread.ID)

	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
