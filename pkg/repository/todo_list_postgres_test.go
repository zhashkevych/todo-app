package repository

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"github.com/zhashkevych/todo-app"
	"testing"
)

func TestTodoListPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		item   todo.TodoList
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO todo_lists").
					WithArgs("title", "description").WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_lists").WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: args{
				userId: 1,
				item: todo.TodoList{
					Title:       "title",
					Description: "description",
				},
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			mock: func() {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO todo_lists").
					WithArgs("", "description").WillReturnRows(rows)

				mock.ExpectRollback()
			},
			input: args{
				userId: 1,
				item: todo.TodoList{
					Title:       "",
					Description: "description",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(tt.input.userId, tt.input.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []todo.TodoList
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1").
					AddRow(2, "title2", "description2").
					AddRow(3, "title3", "description3")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			want: []todo.TodoList{
				{1, "title1", "description1"},
				{2, "title2", "description2"},
				{3, "title3", "description3"},
			},
		},
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1").
					AddRow(2, "title2", "description2").
					AddRow(3, "title3", "description3")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			want: []todo.TodoList{
				{1, "title1", "description1"},
				{2, "title2", "description2"},
				{3, "title3", "description3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll(tt.input.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    todo.TodoList
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: todo.TodoList{1, "title1", "description1"},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 404).WillReturnRows(rows)
			},
			input: args{
				listId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetById(tt.input.userId, tt.input.listId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_lists tl USING users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_lists tl USING users_lists ul WHERE (.+)").
					WithArgs(1, 404).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				listId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input.userId, tt.input.listId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		listId int
		userId int
		input  todo.UpdateListInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				input: todo.UpdateListInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_WithoutDescription",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("new title", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				input: todo.UpdateListInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_WithoutTitle",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("new description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				input: todo.UpdateListInput{
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET FROM users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Update(tt.input.userId, tt.input.listId, tt.input.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
