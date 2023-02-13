package todo

import (
	"errors"
)

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListItem struct {
	Id     int
	UserId string
	ListId int
}

type UpdateListItem struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateListItem) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("add title or description for update values")
	}

	return nil
}

type UpdateItem struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *string `json:"done"`
}

func (i UpdateItem) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("add title, description or done for update values")
	}

	return nil
}
