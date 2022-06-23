package domain

import "time"

type Post struct {
	Id       int64     `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author" validate:"required"`
	Message  string    `json:"message" validate:"required"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int32     `json:"thread"`
	Created  time.Time `json:"created"`
}

type PostInfo struct {
	Post   Post    `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
