package main

import (
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestSelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "title", "updated", "description"})
	expect := []*Item{
		{elemID, "title", "desct", sql.NullString{}},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Title, nil, item.Description)
	}

	mock.
		ExpectQuery("SELECT id, title, updated, description FROM items WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := &ItemRepository{
		DB: db,
	}
	item, err := repo.SelectByID(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.Equalf(t, item, expect[0], "results not match, want %v, have %v", expect[0], item)

	// query error
	mock.
		ExpectQuery("SELECT id, title, updated, description FROM items WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.SelectByID(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "title")

	mock.
		ExpectQuery("SELECT id, title, updated, description FROM items WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = repo.SelectByID(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ItemRepository{
		DB: db,
	}

	title := "title"
	descr := "description"
	testItem := &Item{
		Title:       title,
		Description: descr,
	}

	//ok query
	mock.
		ExpectExec(`INSERT INTO items`).
		WithArgs(title, descr).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.Create(testItem)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if id != 1 {
		t.Errorf("bad id: want %v, have %v", id, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO items`).
		WithArgs(title, descr).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.Create(testItem)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// result error
	mock.
		ExpectExec(`INSERT INTO items`).
		WithArgs(title, descr).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.Create(testItem)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// // last id error
	// mock.
	// 	ExpectExec(`INSERT INTO items`).
	// 	WithArgs(title, descr).
	// 	WillReturnResult(sqlmock.NewResult(0, 0))

	// _, err = repo.Create(testItem)
	// if err == nil {
	// 	t.Errorf("expected error, got nil")
	// 	return
	// }
	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// }

}
