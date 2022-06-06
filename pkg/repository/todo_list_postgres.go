package repository

import (
	"fmt"
	"strings"
	"todo-app"

	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{
		db: db,
	}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin() // запускаем транзакцию
	if err != nil {
		return 0, err
	}

	var id int
	// Первой командой добавляем title и description, полученные от пользователя в таблицу с постами и получаем id нового поста
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Второй командой создаем зависимости пользователя и созданного поста
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit() // Обязательно коммитим транзакцию
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) { // Создаем слайс спизков определенного user`а
	var lists []todo.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) DeleteById(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) UpdateById(userId, listId int, list todo.UpdateListInput) (todo.TodoList, error) {
	// Обновление записей должно происходить как только для поля title, так и только для поля description,
	// так и для обоих полей одновременно. Для это реализуем следующее:
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if list.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *list.Title)
		argId++
	}
	if list.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *list.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	var newList todo.TodoList

	updateListQuery := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.user_id = $%d AND ul.list_id = $%d RETURNING tl.id, tl.title, tl.description",
		todoListsTable, setQuery, usersListsTable, argId, (argId + 1))
	args = append(args, userId, listId)
	err := r.db.QueryRow(updateListQuery, args...).Scan(&newList.Id, &newList.Title, &newList.Description)

	return newList, err
}
