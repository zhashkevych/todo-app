package repository

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type TodoItemRedis struct {
	context     *gin.Context
	redisClient *redis.Client
}

func NewTodoItemRedis(context *gin.Context, redisClient *redis.Client) *TodoItemRedis {
	return &TodoItemRedis{
		context:     context,
		redisClient: redisClient,
	}
}

// Если listId использовать не нужно, передать -1
func (r *TodoItemRedis) HGet(userId, listId, itemId int) (string, error) {
	var val string
	var err error

	if itemId < 0 && listId >= 0 {
		val, err = r.redisClient.HGet(r.context, fmt.Sprintf("user:%d", userId), fmt.Sprintf("items:list%d", listId)).Result()
	} else if itemId >= 0 && listId < 0 {
		val, err = r.redisClient.HGet(r.context, fmt.Sprintf("user:%d", userId), fmt.Sprintf("item:%d", itemId)).Result()
	} else {
		err = errors.New("invalide func HSet")
		return "", err
	}

	return val, err
}

// Если listId использовать не нужно, передать -1
func (r *TodoItemRedis) HSet(userId, listId, itemId int, data string) error {

	//Используем команду конвейер (Pipeline) для одновременного выполнения команд записи в кэш и установление тайм-аута ключа
	pipe := r.redisClient.Pipeline() // создание конвейра

	if itemId < 0 && listId >= 0 {
		pipe.HSetNX(r.context, fmt.Sprintf("user:%d", userId), fmt.Sprintf("items:list%d", listId), data) // Кешируем lists в Redis
	} else if itemId >= 0 && listId < 0 {
		pipe.HSetNX(r.context, fmt.Sprintf("user:%d", userId), fmt.Sprintf("item:%d", itemId), data)
	} else {
		err := errors.New("invalide func HSet")
		return err
	}
	pipe.Expire(r.context, fmt.Sprintf("user:%d", userId), duration) // Устанавливаем тайм-айт для ключа
	_, err := pipe.Exec(r.context)                                   // Выполняем команды конвейера

	return err
}

func (r *TodoItemRedis) HDelete(userId, listId int) error {
	err := r.redisClient.HDel(r.context, fmt.Sprintf("user:%d", userId), fmt.Sprintf("items:list%d", listId)).Err()
	return err
}

func (r *TodoItemRedis) Delete(userId int) error {
	err := r.redisClient.Del(r.context, fmt.Sprintf("user:%d", userId)).Err()
	return err
}
