package repository

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"github.com/zhashkevych/todo-app"
	"testing"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   todo.TodoItem
	}
	type mockBehavior func(args args, id int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			want: 2,
			mock: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "description",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed 2nd Insert",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "title",
					Description: "description",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).
					WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.Create(tt.input.listId, tt.input.item)
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

func TestTodoItemPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []todo.TodoItem
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true).
					AddRow(2, "title2", "description2", false).
					AddRow(3, "title3", "description3", false)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: []todo.TodoItem{
				{1, "title1", "description1", true},
				{2, "title2", "description2", false},
				{3, "title3", "description3", false},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
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

			got, err := r.GetAll(tt.input.userId, tt.input.listId)
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

func TestTodoItemPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    todo.TodoItem
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
			want: todo.TodoItem{1, "title1", "description1", true},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(404, 1).WillReturnRows(rows)
			},
			input: args{
				itemId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetById(tt.input.userId, tt.input.itemId)
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

func TestTodoItemPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
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
				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_lists ul WHERE (.+)").
					WithArgs(1, 404).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				itemId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input.userId, tt.input.itemId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
		userId int
		input  todo.UpdateItemInput
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
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", true, 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: todo.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
					Done:        boolPointer(true),
				},
			},
		},
		{
			name: "OK_WithoutDone",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: todo.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_WithoutDoneAndDescription",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: todo.UpdateItemInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Update(tt.input.userId, tt.input.itemId, tt.input.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
