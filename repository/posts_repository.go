package repository

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"github.com/jackc/pgx"
	"net/http"
	"strconv"
)

func (store *DBStore) PutPost(post *models.Post) (uint64, *models.Error) {
	var ID uint64

	insertQuery := `INSERT INTO posts (created, forumid, isedited, message, parentid, authorid, threadid, parents) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, (SELECT parents FROM posts WHERE posts.id = $5)) RETURNING id`
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

func (store *DBStore) GetPostsByThreadID(threadID int64, params models.PostParams) (models.Posts, *models.Error) {
	posts := models.Posts{}

	var curParams []interface{}
	selectStr := ""

	//`SELECT ID, created, forumid, isEdited, message, parentid, authorid, threadid
	//			FROM posts WHERE threadID = $1`

	switch params.Sort {
	case models.Flat:
		curParams = append(curParams, threadID)
		selectStr += `SELECT p.id, p.created, p.forumid, p.isEdited, 
				p.message, p.parentid, p.authorid, p.threadid FROM posts p WHERE p.threadid = $1`
		if params.Since != -1 {
			curParams = append(curParams, params.Since)
			selectStr += ` AND (p.created, p.id) `
			if params.Desc {
				selectStr += "<"
			} else {
				selectStr += ">"
			}
			selectStr += ` (SELECT posts.created, posts.id FROM posts WHERE posts.id=$2)`
		}
		selectStr += ` ORDER BY (p.created, p.id)`
		if params.Desc {
			selectStr += " DESC"
		}
		if params.Limit != -1 {
			selectStr += " LIMIT $"
			selectStr += strconv.Itoa(len(curParams) + 1)
			curParams = append(curParams, params.Limit)
		}
	case models.Tree:
		curParams = append(curParams, threadID)
		selectStr += `SELECT p.id, p.created, p.forumid, p.isEdited, 
				p.message, p.parentid, p.authorid, p.threadid FROM posts p WHERE p.threadid = $1`
		if params.Since != -1 {
			curParams = append(curParams, params.Since)
			selectStr += " AND p.parents "
			if params.Desc {
				selectStr += "<"
			} else {
				selectStr += ">"
			}
			selectStr += ` (SELECT posts.parents FROM posts WHERE posts.id = $2)`
		}
		selectStr += " ORDER BY p.parents"
		if params.Desc {
			selectStr += " DESC"
		}
		if params.Limit != -1 {
			selectStr += " LIMIT $"
			selectStr += strconv.Itoa(len(curParams) + 1)
			curParams = append(curParams, params.Limit)
		}
	case models.ParentTree:
		curParams = append(curParams, threadID)
		selectStr += `SELECT p.id, p.created, p.forumid, p.isEdited, 
				p.message, p.parentid, p.authorid, p.threadid FROM posts p WHERE p.parents[1] IN (
				SELECT posts.id FROM posts WHERE posts.threadid = $1 AND posts.parentid = 0`
		if params.Since != -1 {
			curParams = append(curParams, params.Since)
			selectStr += ` AND posts.id `
			if params.Desc {
				selectStr += "<"
			} else {
				selectStr += ">"
			}
			selectStr += ` (SELECT COALESCE(posts.parents[1], posts.id) FROM posts WHERE posts.id = $2)`
		}
		selectStr += " ORDER BY posts.id"
		if params.Desc {
			selectStr += " DESC"
		}
		if params.Limit != -1 {
			selectStr += " LIMIT $"
			selectStr += strconv.Itoa(len(curParams) + 1)
			curParams = append(curParams, params.Limit)
		}
		selectStr += `) ORDER BY`
		if params.Desc {
			selectStr += ` p.parents[1] DESC,`
		}
		selectStr += ` p.parents`
	}
	selectStr += ";"

	fmt.Println("НЕ ЖОПА", selectStr, curParams)

	rows, err := store.DB.Query(selectStr, curParams...)
	if err != nil {
		fmt.Println("ЖОПА", selectStr, curParams)
		return posts, models.NewError(http.StatusInternalServerError, err.Error())
	}

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Created, &post.ForumID, &post.IsEdited,
			&post.Message, &post.Parent, &post.AuthorID, &post.Thread)
		if err != nil {
			return posts, models.NewError(http.StatusInternalServerError, err.Error())
		}
		posts = append(posts, post)
	}

	rows.Close()

	fmt.Println("ВЫВОД", selectStr, curParams)

	return posts, nil
}
