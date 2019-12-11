package repository

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
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
