package service

import (
	"todo-app"
	"todo-app/pkg/repository"
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

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) DeleteById(userId, listId int) error {
	return s.repo.DeleteById(userId, listId)
}

func (s *TodoListService) UpdateById(userId, listId int, list todo.UpdateListInput) (todo.TodoList, error) {
	if err := list.Validate(); err != nil {
		var res todo.TodoList
		return res, err
	}
	return s.repo.UpdateById(userId, listId, list)
}
