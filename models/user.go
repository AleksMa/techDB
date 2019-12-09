package models

type User struct {
	ID       int64  `json:"-"`
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

type Users []*User

type UpdateUserFields struct {
	About    *string `json:"about"`
	Email    *string `json:"email"`
	Fullname *string `json:"fullname"`
}
