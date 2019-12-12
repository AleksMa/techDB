package useCase

import (
	"fmt"
	"github.com/AleksMa/techDB/models"
)

func (u *useCase) PutUser(user *models.User) (models.Users, *models.Error) {
	fmt.Println(user)

	if err := user.Validate(); err != nil {
		return nil, err
	}

	/*	users, _ := u.repository.GetDupUsers(user)
		if users != nil && len(users) != 0 {
			fmt.Println("DUP: ", users)
			return users, nil
		}*/

	users, _, err := u.repository.PutUser(user)
	return users, err
}

func (u *useCase) GetUserByNickname(nickname string) (models.User, *models.Error) {
	return u.repository.GetUserByNickname(nickname)
}

func (u *useCase) GetUserByID(id int64) (models.User, *models.Error) {
	return u.repository.GetUserByID(id)
}

func (u *useCase) ChangeUser(userUpd *models.UpdateUserFields, nickname string) (models.User, *models.Error) {
	tempUser, err := u.GetUserByNickname(nickname)
	if err != nil {
		return tempUser, err
	}

	if userUpd.Email != nil {
		tempUser.Email = *userUpd.Email
	}
	if userUpd.Fullname != nil {
		tempUser.Fullname = *userUpd.Fullname
	}
	if userUpd.About != nil {
		tempUser.About = *userUpd.About
	}

	err = u.repository.ChangeUser(&tempUser)
	fmt.Println(*userUpd)
	fmt.Println(tempUser)
	return tempUser, err
}

func (u *useCase) GetUsersByForum(slug string, params models.UserParams) (models.Users, *models.Error) {
	forum, e := u.repository.GetForumBySlug(slug)
	if e != nil {
		return nil, e
	}

	users, e := u.repository.GetUsersByForum(forum.ID, params)
	if e != nil {
		return nil, e
	}

	fmt.Println(users)
	for i, _ := range users {
		user, _ := u.repository.GetUserByID(users[i].ID)
		users[i].Nickname = user.Nickname
		users[i].Fullname = user.Fullname
		users[i].About = user.About
		users[i].Email = user.Email
	}
	return users, nil
}
