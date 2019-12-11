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

func (u *useCase) PutVote(vote *models.Vote) (models.Vote, error) {
	fmt.Println(vote)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(vote.Nickname)
	vote.AuthorID = user.ID

	fmt.Println(vote)

	u.repository.PutVote(vote)
	//TODO: error check
	return *vote, nil
}

func (u *useCase) ChangePost(post *models.Post) error {
	fmt.Println(post)
	//TODO: contains check
	u.repository.ChangePost(post)
	//TODO: error check
	return nil
}

func (u *useCase) GetPostFull(id int64) (models.PostFull, error) {
	var postFull models.PostFull
	var err error
	postFull.Post, err = u.repository.GetPost(id)
	fmt.Println(postFull)
	fmt.Println(err)

	postFull.Author, _ = u.repository.GetUserByID(postFull.Post.AuthorID)
	postFull.Thread, _ = u.repository.GetThreadByID(postFull.Post.Thread)
	postFull.Forum, _ = u.GetForumByID(postFull.Post.ForumID)

	return postFull, nil
}

func (u *useCase) PutVoteWithSlug(vote *models.Vote, slug string) (models.Vote, error) {
	fmt.Println(vote)
	//TODO: contains check
	user, _ := u.repository.GetUserByNickname(vote.Nickname)
	vote.AuthorID = user.ID
	thread, _ := u.repository.GetThreadBySlug(slug)
	vote.ThreadID = thread.ID
	fmt.Println(vote)

	u.repository.PutVote(vote)
	//TODO: error check
	return *vote, nil
}
