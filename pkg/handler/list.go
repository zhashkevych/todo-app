package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo-app"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// @Summary Create todo List
// @Security ApiKeyAuth
// @Tags lists
// @Description create todo List
// @ID create-list
// @Accept json
// @Produce json
// @Param input body todo.TodoList true "List info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c) // Определяем ID юзера по токену
	if err != nil {
		return
	}

	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil { // парсим тело запроса в структуру List
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input) // Создаем список в базе данных
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Удалим список lists из кэша redis
	err = h.services.TodoListCach.HDelete(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{ // Отвечаем ОК, id list
		"id": id,
	})
}

type getAllListsResponce struct { // Структура для использования в ответе
	Data []todo.TodoList `json:"data"`
}

// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all lists
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllListsResponce
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	lists := make([]todo.TodoList, 0)

	userId, err := getUserId(c) // Определяем ID юзера по токену
	if err != nil {
		return
	}

	// Проверяем существует ли ключ "lists_userId" в кэше redis
	val, err := h.services.TodoListCach.HGet(userId, -1)
	if err == redis.Nil { // Если ключа не существует, вытаскиваем данные из postgres и кэшируем в redis

		logrus.Print("Request to Postgres")

		lists, err = h.services.TodoList.GetAll(userId) // вытаскиваем списки из БД для определенного пользователя
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := json.Marshal(lists) // декодируем JSON в слайз байт для дальнейшей записи в redis
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Добавим list в кэш Redis.
		err = h.services.TodoListCach.HSet(userId, -1, string(data))
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	} else if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else { // Если в redis есть ключ...
		logrus.Print("Request to Redis")
		json.Unmarshal([]byte(val), &lists) // забираем от туда данные и отправляем
	}
	c.JSON(http.StatusOK, getAllListsResponce{
		Data: lists,
	})

}

func (h *Handler) getListById(c *gin.Context) {
	var list todo.TodoList

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "ivalid user id")
		return
	}

	id, err := strconv.Atoi(c.Param("id")) // парсим URL, определяем id списка
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid type list id")
		return
	}

	// Проверяем существует ли ключ "hlists_userId" с полем list:id в хэш-таблице redis
	val, err := h.services.TodoListCach.HGet(userId, id)

	if err == redis.Nil { // Если ключа не существует, вытаскиваем данные из postgres и кэшируем в redis

		logrus.Print("Request to Postgres")

		list, err = h.services.TodoList.GetById(userId, id) // вытаскиваем из БД список по id списка и пользователя
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := json.Marshal(list) // декодируем list в слайз байт для дальнейшей записи в redis
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Добавим list в кэш Redis.
		err = h.services.TodoListCach.HSet(userId, id, string(data))
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	} else if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else { // Если в redis есть ключ...
		logrus.Print("Request to Redis")
		json.Unmarshal([]byte(val), &list) // забираем от туда данные и отправляем
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "ivalid user id")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid type list id")
		return
	}

	var input todo.UpdateListInput
	if err := c.BindJSON(&input); err != nil { // парсим тело запроса в структуру List
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	list, err := h.services.TodoList.UpdateById(userId, id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Удаляем все данные из кэша Redis, т.к. изменения могли коснуться любого поля ключа user:userId
	err = h.services.TodoListCach.Delete(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary Delete todo List
// @Security ApiKeyAuth
// @Tags lists
// @Descriprion gelete list by id
// @ID delete-list
// @Accept json
// @Produce json
// @Param id path int true "List Id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "ivalid user id")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid type list id")
		return
	}

	err = h.services.TodoList.DeleteById(userId, id) // Удаляем из таблицы Списков и связывающей таблицы список по id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Удаляем все данные из кэша Redis, т.к. изменения могли коснуться любого поля ключа user:userId
	err = h.services.TodoListCach.Delete(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Ok": fmt.Sprintf("deleted list by id: %d", id),
	})
}
