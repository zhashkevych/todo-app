package repository

import (
	"todo-app"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	DeleteById(userId, listId int) error
	UpdateById(userId, listId int, list todo.UpdateListInput) (todo.TodoList, error)
}

type TodoItem interface {
	Create(listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}

type TodoListCach interface {
	HGet(userId, listId int) (string, error)
	HSet(userId, listId int, data string) error
	HDelete(userId int) error
	Delete(userId int) error
}

type TodoItemCach interface {
	HGet(userId, listId, itemId int) (string, error)
	HSet(userId, listId, itemId int, data string) error
	HDelete(userId, listId int) error
	Delete(userId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
	TodoListCach
	TodoItemCach
}

func NewRepository(db *sqlx.DB, context *gin.Context, redisClient *redis.Client) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
		TodoListCach:  NewTodoListRedis(context, redisClient),
		TodoItemCach:  NewTodoItemRedis(context, redisClient),
	}

}
