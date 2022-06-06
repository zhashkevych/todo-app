package service

import (
	"todo-app/pkg/repository"
)

type TodoListServiceCach struct {
	repo repository.TodoListCach
}

func NewTodoListServiceCach(repo repository.TodoListCach) *TodoListServiceCach {
	return &TodoListServiceCach{repo: repo}
}

// Если listId использовать не нужно, передать -1
func (s *TodoListServiceCach) HGet(userId, listId int) (string, error) {
	return s.repo.HGet(userId, listId)
}

// Если listId использовать не нужно, передать -1
func (s *TodoListServiceCach) HSet(userId, listId int, data string) error {
	return s.repo.HSet(userId, listId, data)
}

func (s *TodoListServiceCach) HDelete(userId int) error {
	return s.repo.HDelete(userId)
}

func (s *TodoListServiceCach) Delete(userId int) error {
	return s.repo.Delete(userId)
}
