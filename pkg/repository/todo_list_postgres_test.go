package repository

/*func TestTodoListPostgres_Create(t *testing.T) {
	// Создаем зависимости
	db, mock, err := sqlmock.Newx() // mock объекта БД
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		listId int
		item todo.TodoItem
	}
	type mockBehavior func(args args, id int)

	testTable := []struct {
		name string
		mockBehavior mockBehavior
		args args
		id int
		wantErr bool // флаг ожидания ошибки
	}
}*/
