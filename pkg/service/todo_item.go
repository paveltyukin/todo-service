package service

import (
	"todo"
	"todo/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func (s *TodoItemService) Update(userId, itemId int, item todo.UpdateItem) error {
	if err := item.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, itemId, item)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, input todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, input)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAll(listId, userId)
}

func (s *TodoItemService) GetById(userId int, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}
