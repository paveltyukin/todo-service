package service

import (
	"todo"
	"todo/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId int, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) Delete(userId int, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, item todo.UpdateListItem) error {
	if err := item.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, listId, item)
}
