package repository

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var duration time.Duration = 3600 * time.Second // тайм-айт ключа в кэше Redis

type TodoListRedis struct {
	context     *gin.Context
	redisClient *redis.Client
}

func NewTodoListRedis(context *gin.Context, redisClient *redis.Client) *TodoListRedis {
	return &TodoListRedis{
		context:     context,
		redisClient: redisClient,
	}
}

// Если listId использовать не нужно, передать -1
func (r *TodoListRedis) HGet(userId, listId int) (string, error) {
	if listId < 0 {
		val, err := r.redisClient.HGet(r.context, fmt.Sprintf("user:%d", userId), "lists").Result() // Запрос значения из кэша по ключу userId и полю lists
		return val, err
	} else {
		val, err := r.redisClient.HGet(r.context, fmt.Sprintf("user:%d", userId), fmt.Sprintf("list:%d", listId)).Result()
		return val, err
	}
}

// Если listId использовать не нужно, передать -1
func (r *TodoListRedis) HSet(userId, listId int, data string) error {

	//Используем команду конвейер (Pipeline) для одновременного выполнения команд записи в кэш и установление тайм-аута ключа
	pipe := r.redisClient.Pipeline() // создание конвейра

	if listId < 0 {
		pipe.HSetNX(r.context, fmt.Sprintf("user:%d", userId), "lists", data) // Кешируем lists в Redis
	} else {
		pipe.HSetNX(r.context, fmt.Sprintf("user:%d", userId), fmt.Sprintf("list:%d", listId), data) // Кешируем list в Redis
	}
	pipe.Expire(r.context, fmt.Sprintf("user:%d", userId), duration) // Устанавливаем тайм-айт для ключа
	_, err := pipe.Exec(r.context)                                   // Выполняем команды конвейера

	return err
}

func (r *TodoListRedis) HDelete(userId int) error {
	err := r.redisClient.HDel(r.context, fmt.Sprintf("user:%d", userId), "lists").Err()
	return err
}

func (r *TodoListRedis) Delete(userId int) error {
	err := r.redisClient.Del(r.context, fmt.Sprintf("user:%d", userId)).Err()
	return err
}
