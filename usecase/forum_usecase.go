package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
	"net/http"
)

func (u *useCase) PutForum(forum *models.Forum) (models.Forum, *models.Error) {
	fmt.Println(forum)
	dupForum, err := u.GetForumBySlug(forum.Slug)
	if err == nil || err.Code != http.StatusNotFound {
		fmt.Println("DUP: ", dupForum)
		return dupForum, models.NewError(http.StatusConflict, "forum already created")
	}

	user, err := u.GetUserByNickname(forum.Owner)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return *forum, models.NewError(http.StatusNotFound, "No user found: "+err.Message)
		}
		return *forum, models.NewError(http.StatusInternalServerError, err.Message)
	}

	forum.OwnerID = user.ID
	forum.Owner = user.Nickname
	_, err = u.repository.PutForum(forum)

	if err != nil {
		return *forum, err
	}
	return *forum, nil
}

func (u *useCase) GetForumBySlug(slug string) (models.Forum, *models.Error) {
	forum, err := u.repository.GetForumBySlug(slug)
	if err != nil {
		return forum, err
	}

	fmt.Println(forum)
	owner, err := u.repository.GetUserByID(forum.OwnerID)
	if err != nil {
		return forum, err
	}

	forum.Owner = owner.Nickname
	return forum, nil
}

func (u *useCase) GetForumByID(id int64) (models.Forum, *models.Error) {
	forum, err := u.repository.GetForumByID(id)
	if err != nil {
		return forum, err
	}

	fmt.Println(forum)
	owner, err := u.repository.GetUserByID(forum.OwnerID)
	if err != nil {
		return forum, err
	}

	forum.Owner = owner.Nickname
	return forum, nil
}
