package repository

import (
	"github.com/jmoiron/sqlx"
	"todo"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetById(userId int, listId int) (todo.TodoList, error)
	GetAll(userId int) ([]todo.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, item todo.UpdateListItem) error
}

type TodoItem interface {
	Create(listId int, input todo.TodoItem) (int, error)
	GetAll(listId, userId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, item todo.UpdateItem) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
