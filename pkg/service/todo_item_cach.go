package service

import (
	"todo-app/pkg/repository"
)

type TodoItemServiceCach struct {
	repo repository.TodoItemCach
}

func NewTodoItemServiceCach(repo repository.TodoItemCach) *TodoItemServiceCach {
	return &TodoItemServiceCach{repo: repo}
}

// Если listId или itemId использовать не нужно, передать -1
func (s *TodoItemServiceCach) HGet(userId, listId, itemId int) (string, error) {
	return s.repo.HGet(userId, listId, itemId)
}

// Если listId или itemId использовать не нужно, передать -1
func (s *TodoItemServiceCach) HSet(userId, listId, itemId int, data string) error {
	return s.repo.HSet(userId, listId, itemId, data)
}

func (s *TodoItemServiceCach) HDelete(userId, listId int) error {
	return s.repo.HDelete(userId, listId)
}

func (s *TodoItemServiceCach) Delete(userId int) error {
	return s.repo.Delete(userId)
}
