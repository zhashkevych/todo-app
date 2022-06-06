package service

import (
	"todo-app"
	"todo-app/pkg/repository"
)

// Создание mock интерфейсов в папке mocks
//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	DeleteById(userId, listId int) error
	UpdateById(userId, listId int, list todo.UpdateListInput) (todo.TodoList, error)
}

type TodoItem interface {
	Create(userId, listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}

type TodoListCach interface {
	// Если listId использовать не нужно, передать -1
	HGet(userId, listId int) (string, error)
	// Если listId использовать не нужно, передать -1
	HSet(userId, listId int, data string) error
	HDelete(userId int) error
	Delete(userId int) error
}

type TodoItemCach interface {
	// Если listId или itemId использовать не нужно, передать -1 (что то одно)
	HGet(userId, listId, itemId int) (string, error)
	// Если listId или itemId использовать не нужно, передать -1 (что то одно)
	HSet(userId, listId, itemId int, data string) error
	HDelete(userId, listId int) error
	Delete(userId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
	TodoListCach
	TodoItemCach
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
		TodoListCach:  NewTodoListServiceCach(repos.TodoListCach),
		TodoItemCach:  NewTodoItemServiceCach(repos.TodoItemCach),
	}
}
