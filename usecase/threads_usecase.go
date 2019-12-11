package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"net/http"
)

func (u *useCase) PutThread(thread *models.Thread) (models.Thread, *models.Error) {
	fmt.Println(thread)
	if thread.Slug != "" {
		dupForum, err := u.GetThreadBySlug(thread.Slug)
		if err == nil || err.Code != http.StatusNotFound {
			fmt.Println("DUP: ", dupForum)
			return dupForum, models.NewError(http.StatusConflict, "forum already created")
		}
	}

	user, err := u.GetUserByNickname(thread.Author)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return *thread, models.NewError(http.StatusNotFound, "No user found: "+err.Message)
		}
		return *thread, models.NewError(http.StatusInternalServerError, err.Message)
	}

	forum, err := u.GetForumBySlug(thread.Forum)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return *thread, models.NewError(http.StatusNotFound, "No forum found: "+err.Message)
		}
		return *thread, models.NewError(http.StatusInternalServerError, err.Message)
	}
	thread.AuthorID = user.ID
	thread.ForumID = forum.ID
	thread.Forum = forum.Slug

	id, err := u.repository.PutThread(thread)
	if err != nil {
		return *thread, err
	}
	thread.ID = int64(id)
	//thread.Created = thread.Created.Add(time.Duration(-3) * time.Hour)
	return *thread, nil
}

func (u *useCase) GetThreadBySlug(slug string) (models.Thread, *models.Error) {
	thread, err := u.repository.GetThreadBySlug(slug)
	if err != nil {
		return thread, err
	}

	fmt.Println(thread)

	owner, err := u.GetUserByID(thread.AuthorID)
	if err != nil {
		return thread, err
	}
	thread.Author = owner.Nickname

	forum, err := u.repository.GetForumByID(thread.ForumID)
	if err != nil {
		return thread, err
	}

	thread.Forum = forum.Slug
	return thread, nil
}

func (u *useCase) GetThreadByID(id int64) (models.Thread, *models.Error) {
	thread, err := u.repository.GetThreadByID(id)
	if err != nil {
		return thread, err
	}

	fmt.Println(thread)

	owner, err := u.GetUserByID(thread.AuthorID)
	if err != nil {
		return thread, err
	}
	thread.Author = owner.Nickname

	forum, err := u.repository.GetForumByID(thread.ForumID)
	if err != nil {
		return thread, err
	}

	thread.Forum = forum.Slug
	return thread, nil
}

func (u *useCase) GetThreadsByForum(slug string, params models.ThreadParams) (models.Threads, *models.Error) {
	forum, err := u.GetForumBySlug(slug)
	if err != nil {
		return nil, err
	}

	threads, _ := u.repository.GetThreadsByForum(forum.ID, params)
	for i, _ := range threads {
		threads[i].Forum = forum.Slug
		user, _ := u.repository.GetUserByID(threads[i].AuthorID)
		threads[i].Author = user.Nickname
	}
	return threads, nil
}

func (u *useCase) UpdateThreadWithID(thread *models.Thread) (models.Thread, error) {
	fmt.Println(thread)
	//TODO: contains check
	u.repository.UpdateThreadWithID(thread)
	//TODO: error check
	return *thread, nil
}

func (u *useCase) UpdateThreadWithSlug(thread *models.Thread) (models.Thread, error) {
	fmt.Println(thread)
	//TODO: contains check
	u.repository.UpdateThreadWithSlug(thread)
	//TODO: error check
	return *thread, nil
}
