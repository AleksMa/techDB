package repository

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/jackc/pgx"
	"net/http"
)

func (store *DBStore) PutPost(post *models.Post) (uint64, *models.Error) {
	fmt.Println(post)
	var ID uint64

	insertQuery := `INSERT INTO posts (created, forumid, isedited, message, parentid, authorid, threadid) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	rows := store.DB.QueryRow(insertQuery,
		post.Created, post.ForumID, post.IsEdited, post.Message, post.Parent, post.AuthorID, post.Thread)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return ID, nil
}

func (store *DBStore) GetPost(id int64) (models.Post, *models.Error) {
	post := &models.Post{}

	selectStr := "SELECT ID, created, forumid, isEdited, message, parentid, authorid, threadid FROM posts WHERE id = $1"
	row := store.DB.QueryRow(selectStr, id)

	err := row.Scan(&post.ID, &post.Created, &post.ForumID, &post.IsEdited,
		&post.Message, &post.Parent, &post.AuthorID, &post.Thread)

	if err != nil {
		if err == pgx.ErrNoRows {
			return *post, models.NewError(http.StatusNotFound, err.Error())
		}
		return *post, models.NewError(http.StatusInternalServerError, err.Error())
	}

	return *post, nil
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

func (store *DBStore) PutVote(vote *models.Vote) (uint64, *models.Error) {
	tx, _ := store.DB.Begin()
	defer tx.Rollback()

	fmt.Println(vote)
	var ID uint64
	positive := vote.Voice == 1

	insertQuery := `INSERT INTO votes (voice, threadid, authorid) VALUES ($1, $2, $3) RETURNING id`
	rows := tx.QueryRow(insertQuery, positive, vote.ThreadID, vote.AuthorID)

	err := rows.Scan(&ID)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, models.NewError(http.StatusInternalServerError, "Can not put vote: "+err.Error())
	}

	updateQuery := `UPDATE threads SET vote=vote+$1 WHERE id=$2`
	_, err = tx.Exec(updateQuery,
		vote.Voice, vote.ThreadID)

	if err != nil {
		fmt.Println("ERR: ", err.Error())
		if pgerr, ok := err.(pgx.PgError); ok && pgerr.Code == "23505" {
			return 0, models.NewError(http.StatusNotFound, err.Error())
		}
		return 0, models.NewError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()

	return ID, nil
}

func (store *DBStore) ChangePost(post *models.Post) *models.Error {
	fmt.Println(post)

	insertQuery := `UPDATE posts SET message=$1, isedited=$2 WHERE id=$3`
	_, err := store.DB.Exec(insertQuery, post.Message, true, post.ID)

	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (store *DBStore) UpdateVote(vote *models.Vote) (int, *models.Error) {
	tx, _ := store.DB.Begin()
	defer tx.Rollback()

	tempVote := -1
	var tempVoice bool

	fmt.Println(vote)

	selectStr := "SELECT voice FROM votes WHERE authorid=$1 AND threadid=$2"
	row := store.DB.QueryRow(selectStr, vote.AuthorID, vote.ThreadID)

	err := row.Scan(&tempVoice)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return 0, models.NewError(http.StatusConflict, err.Error())
		}
		return 0, models.NewError(http.StatusInternalServerError, err.Error())
	}

	if tempVoice {
		tempVote = 1
	}

	if tempVote == vote.Voice {
		tempVote = 0
	} else {
		tempVote = vote.Voice * 2
	}

	updateQuery := `UPDATE votes SET voice=$1 WHERE authorid=$2 AND threadid=$3`
	_, err = tx.Exec(updateQuery,
		vote.Voice == 1, vote.AuthorID, vote.ThreadID)

	if err != nil {
		return tempVote, models.NewError(http.StatusConflict, err.Error())
	}

	updateQuery = `UPDATE threads SET vote=vote+$1 WHERE id=$2`
	_, err = tx.Exec(updateQuery,
		tempVote, vote.ThreadID)

	if err != nil {
		fmt.Println("ERR: ", err.Error())
		if pgerr, ok := err.(pgx.PgError); ok && pgerr.Code == "23505" {
			return tempVote, models.NewError(http.StatusNotFound, err.Error())
		}
		return tempVote, models.NewError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()

	return tempVote, nil
}
