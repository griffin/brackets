package env

import (
	"time"
)

const (
	createPost = ""
	getPost = ""
	getPosts = ""
	updatePost = ""
	deletePost = ""
)

type postDatastore interface {
	CreatePost(team Post) (*Post, error)
	GetPost(selector string) (*Post, error)
	GetPosts(selector string) ([]*Post, error)
	UpdatePost(team Team) error
	DeletePost(selector string) error
}


type Post struct {
	Selectable

	id     uint
	Title  string
	Author *User
	Path   string
	Posted time.Time
}


func (d *db) CreatePost(team Post) (*Post, error){
	return nil, nil
}

func (d *db) GetPost(selector string) (*Post, error) {
	return nil, nil
}

func (d *db) GetPosts(teamSelector string) ([]*Post, error){
	return nil, nil
}

func (d *db) UpdatePost(team Team) error {
	return nil
}

func (d *db) DeletePost(selector string) error {
	return nil
}