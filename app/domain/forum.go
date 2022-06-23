package domain

import (
	"time"

	"github.com/Kostich31/techpark_db/app/tools"
)

type Forum struct {
	Title   string `json:"title" validate:"required"`
	User    string `json:"user" validate:"required"`
	Slug    string `json:"slug" validate:"required"`
	Posts   int64  `json:"posts"`
	Threads int64  `json:"threads"`
}

type Thread struct {
	Id      int32     `json:"id"`
	Title   string    `json:"title" validate:"required"`
	Author  string    `json:"author" validate:"required"`
	Forum   string    `json:"forum"`
	Message string    `json:"message" validate:"required"`
	Votes   int32     `json:"votes"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created"`
}

type Vote struct {
	NickName string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type ForumRepository interface {
	AddForum(forum Forum) (Forum, error)
	GetForumBySlug(slug string) (Forum, error)
	GetDetailsForum(slug string) (Forum, error)
	AddThread(thread Thread) (Thread, error)
	GetUsersForum(slug string, filter tools.FilterUser) ([]User, error)
	GetForumThreads(slug string, filter tools.FilterThread) ([]Thread, error)
}

type ForumUseCase interface {
	CreateForum(forum Forum) (Forum, *CustomError)
	GetDetailsForum(slug string) (Forum, *CustomError)
	CreateThread(thread Thread) (Thread, *CustomError)
	GetUsersForum(slug string, filter tools.FilterUser) ([]User, *CustomError)
	GetForumThreads(slug string, filter tools.FilterThread) ([]Thread, *CustomError)
}

type ThreadRepository interface {
	CreatePosts(threadId int, threadForum string, post []Post) ([]Post, error)
	GetThreadBySlug(slug string) (Thread, error)
	GetThreadById(id int) (Thread, error)
	GetThreadBySlugOrId(slugOrId string) (Thread, error)
	CreateVoteBySlugOrId(slugOrId string, vote Vote) error
	UpdateVoteBySlugOrId(slugOrId string, vote Vote) error
	GetPostById(id int) (Post, error)
	UpdatePost(id int, post Post) (Post, error)
	GetPostsFlatSlugOrId(slugOrId string, posts tools.FilterPosts) ([]*Post, error)
	GetPostsTreeSlugOrId(slugOrId string, posts tools.FilterPosts) ([]*Post, error)
	GetPostsParentTreeSlugOrId(slugOrId string, posts tools.FilterPosts) ([]*Post, error)
	UpdateThread(slugOrId string, thread Thread) (Thread, error)
}

type ThreadUseCase interface {
	CreatePosts(slugOrId string, post []Post) ([]Post, *CustomError)
	CreateVote(slugOrId string, vote Vote) (Thread, *CustomError)
	GetThreadDetails(slugOrId string) (Thread, *CustomError)
	GetPosts(slugOrId string, filter tools.FilterPosts) ([]*Post, *CustomError)
	GetPost(id string, filter tools.FilterOnePost) (PostInfo, *CustomError)
	UpdateThread(slugOrId string, thread Thread) (Thread, *CustomError)
	UpdatePost(id string, post Post) (Post, *CustomError)
}
