// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateTodoInput struct {
	Content string `json:"content"`
}

type Todo struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
