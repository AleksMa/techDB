package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"net/http"
)

func (u *useCase) PutPost(post *models.Post) (*models.Post, *models.Error) {
	fmt.Println(post)
	user, err := u.GetUserByNickname(post.Author)
	if err != nil {
		return post, err
	}
	thread, _ := u.GetThreadByID(post.Thread)
	post.Forum = thread.Forum
	post.ForumID = thread.ForumID
	post.AuthorID = user.ID
	//fmt.Println(post, post.Forum)

	if post.Parent != 0 {
		_, err := u.repository.GetPost(post.Parent)
		if err != nil {
			return post, models.NewError(http.StatusConflict, err.Message)
		}
	}

	fmt.Println(post)

	id, err := u.repository.PutPost(post)
	if err != nil {
		return post, err
	}
	post.ID = int64(id)
	return post, nil
}

func (u *useCase) PutPostWithSlug(post *models.Post, threadSlug string) (*models.Post, *models.Error) {
	fmt.Println(post)
	user, err := u.GetUserByNickname(post.Author)
	if err != nil {
		return post, err
	}
	thread, err := u.GetThreadBySlug(threadSlug)
	fmt.Println("SLUG: ", threadSlug)
	if err != nil {
		return post, models.NewError(http.StatusConflict, err.Message)
	}
	post.Thread = thread.ID
	post.Forum = thread.Forum
	post.ForumID = thread.ForumID
	post.AuthorID = user.ID

	if post.Parent != 0 {
		_, err := u.repository.GetPost(post.Parent)
		if err != nil {
			return post, models.NewError(http.StatusConflict, err.Message)
		}
	}

	fmt.Println(post)

	id, err := u.repository.PutPost(post)
	if err != nil {
		return post, err
	}
	post.ID = int64(id)
	return post, nil
}

func (u *useCase) PutVote(vote *models.Vote) (models.Thread, *models.Error) {
	var thread models.Thread
	var voice int
	if vote.Voice != -1 && vote.Voice != 1 {
		return thread, models.NewError(http.StatusInternalServerError, "Unexpected voice")
	}
	fmt.Println(vote)
	user, err := u.GetUserByNickname(vote.Nickname)
	if err != nil {
		return thread, models.NewError(http.StatusInternalServerError, "No such user")
	}
	vote.AuthorID = user.ID

	thread, err = u.GetThreadByID(vote.ThreadID)
	if err != nil {
		return thread, models.NewError(http.StatusNotFound, "No such thread")
	}

	voice, err = u.repository.UpdateVote(vote)
	if err == nil {
		thread.Votes += int32(voice)
		return thread, nil
	}

	_, err = u.repository.PutVote(vote)
	if err != nil {
		return thread, err
	}

	thread.Votes += int32(vote.Voice)
	fmt.Println(vote)

	return thread, nil
}

func (u *useCase) PutVoteWithSlug(vote *models.Vote, slug string) (models.Thread, *models.Error) {
	var thread models.Thread
	var voice int

	if vote.Voice != -1 && vote.Voice != 1 {
		return thread, models.NewError(http.StatusInternalServerError, "Unexpected voice")
	}
	fmt.Println(vote)
	user, err := u.GetUserByNickname(vote.Nickname)
	if err != nil {
		return thread, models.NewError(http.StatusInternalServerError, "No such user")
	}
	vote.AuthorID = user.ID

	thread, err = u.GetThreadBySlug(slug)
	if err != nil {
		return thread, models.NewError(http.StatusNotFound, "No such thread")
	}
	vote.ThreadID = thread.ID

	voice, err = u.repository.UpdateVote(vote)
	if err == nil {
		thread.Votes += int32(voice)
		return thread, nil
	}

	_, err = u.repository.PutVote(vote)
	if err != nil {
		return thread, err
	}

	thread.Votes += int32(vote.Voice)
	fmt.Println(vote)

	return thread, nil
}

func (u *useCase) GetPostFull(id int64, fields []string) (models.PostFull, *models.Error) {
	var postFull models.PostFull
	//var err error
	post, err := u.repository.GetPost(id)
	if err != nil {
		return postFull, err
	}
	postFull.Post = &post
	fmt.Println(postFull)
	fmt.Println(err)

	author, _ := u.GetUserByID(postFull.Post.AuthorID)
	forum, _ := u.GetForumByID(postFull.Post.ForumID)

	postFull.Post.Author = author.Nickname
	postFull.Post.Forum = forum.Slug

	for _, field := range fields {
		if field == "user" {
			postFull.Author = &author
		}

		if field == "forum" {
			postFull.Forum = &forum
		}

		if field == "thread" {
			thread, _ := u.GetThreadByID(postFull.Post.Thread)
			postFull.Thread = &thread
		}
	}

	return postFull, nil
}

func (u *useCase) GetPostsByThreadID(id int64) (models.Posts, error) {
	thread, _ := u.repository.GetThreadByID(id)

	posts, _ := u.repository.GetPostsByThreadID(thread.ID)
	for i, _ := range posts {
		posts[i].Thread = thread.ID
		user, _ := u.repository.GetUserByID(posts[i].AuthorID)
		posts[i].Author = user.Nickname
	}
	return posts, nil
}

func (u *useCase) GetPostsByThreadSlug(slug string) (models.Posts, error) {
	thread, _ := u.repository.GetThreadBySlug(slug)

	posts, _ := u.repository.GetPostsByThreadID(thread.ID)
	for i, _ := range posts {
		posts[i].Thread = thread.ID
		user, _ := u.repository.GetUserByID(posts[i].AuthorID)
		posts[i].Author = user.Nickname
	}
	return posts, nil
}

func (u *useCase) ChangePost(post *models.Post) error {
	fmt.Println(post)
	//TODO: contains check
	u.repository.ChangePost(post)
	//TODO: error check
	return nil
}
