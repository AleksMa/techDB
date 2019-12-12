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
	thread, err := u.GetThreadByID(post.Thread)
	if err != nil {
		return post, err
	}
	post.Forum = thread.Forum
	post.ForumID = thread.ForumID
	post.AuthorID = user.ID

	if post.Parent != 0 {
		parentPost, err := u.GetPostByID(post.Parent)
		if err != nil {
			return post, models.NewError(http.StatusConflict, err.Message)
		}
		if post.Thread != parentPost.Thread {
			return post, models.NewError(http.StatusConflict, "Cross-thread exception")
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
	if err != nil {
		return post, err
	}
	fmt.Println("SLUG: ", threadSlug)
	if err != nil {
		return post, models.NewError(http.StatusConflict, err.Message)
	}
	post.Thread = thread.ID
	post.Forum = thread.Forum
	post.ForumID = thread.ForumID
	post.AuthorID = user.ID

	fmt.Println("POST_PARENT: ", post.Parent)

	if post.Parent != 0 {
		parentPost, err := u.GetPostByID(post.Parent)
		if err != nil {
			return post, models.NewError(http.StatusConflict, err.Message)
		}
		if post.Thread != parentPost.Thread {
			return post, models.NewError(http.StatusConflict, "Cross-thread exception")
		}
	}

	fmt.Printf("%#v", post)

	id, err := u.repository.PutPost(post)
	if err != nil {
		return post, err
	}
	post.ID = int64(id)
	return post, nil
}

func (u *useCase) GetPostByID(id int64) (models.Post, *models.Error) {

	post, err := u.repository.GetPost(id)
	if err != nil {
		return post, err
	}

	fmt.Println(post)

	owner, err := u.GetUserByID(post.AuthorID)
	if err != nil {
		return post, err
	}
	post.Author = owner.Nickname
	fmt.Println(owner)

	forum, err := u.repository.GetForumByID(post.ForumID)
	if err != nil {
		return post, err
	}

	post.Forum = forum.Slug

	fmt.Println(post)
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

func (u *useCase) ChangePost(post *models.Post) (models.Post, *models.Error) {
	tempPost, err := u.GetPostByID(post.ID)
	if err != nil {
		return tempPost, err
	}

	if post.Message == "" || post.Message == tempPost.Message {
		return tempPost, nil
	}

	err = u.repository.ChangePost(post)
	if err != nil {
		return tempPost, err
	}
	tempPost.IsEdited = true
	tempPost.Message = post.Message

	fmt.Println(tempPost)
	return tempPost, err
}
