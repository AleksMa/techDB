package models

import "time"

type Thread struct {
	Author  string     `json:"author"`
	Created *time.Time `json:"created"`
	Forum   string     `json:"forum"`
	ID      int64      `json:"id"`
	Message string     `json:"message"`
	Slug    *string    `json:"slug"`
	Title   string     `json:"title"`
	Votes   int32      `json:"votes"`
}

type Threads []*Thread

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}
