package repository

import (
	"errors"
	"log"
	"testing"

	todo "github.com/Mobo140/projects/to_do_list"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   todo.TodoItem
	}
	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO list_items").WithArgs(id, args.listId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields", //Rollback при первом insert
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Insert Error", //Rollback при втором insert
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO list_items").WithArgs(id, args.listId).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.listId, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err) //Тест упадет если ошибка окажется пустой.
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestTodoItemPostgres_GetAll(t *testing.T) {

}
