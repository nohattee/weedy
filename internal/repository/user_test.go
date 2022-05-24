package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"m0n3y-api/internal/entity"
	"m0n3y-api/internal/helper"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_List(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := UserRepo{DB: db}

	testCases := []helper.TestCase{
		{
			Name: "happy case",
			Req:  &ListParams{},
			Setup: func(ctx context.Context) {
				sqlFieldMap, err := helper.GetSQLFieldMapFromEntity(&entity.User{})
				if err != nil {
					t.Fatalf("an error '%s' was not expected when getting sqlFieldMap", err)
				}
				mock.ExpectQuery("SELECT (.+) users").
					WillReturnRows(
						sqlmock.NewRows(sqlFieldMap.Keys).
							AddRow("student-01", "FirstName 01", "LastName 01", "student-01@example.com", time.Now(), time.Now(), nil),
					)
			},
		},
		{
			Name:        "error: no rows",
			Req:         &ListParams{},
			ExpectedErr: sql.ErrNoRows,
			Setup: func(ctx context.Context) {
				mock.ExpectQuery("SELECT (.+) users").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name:        "error: bad connection",
			Req:         &ListParams{},
			ExpectedErr: sqlmock.ErrCancelled,
			Setup: func(ctx context.Context) {
				sqlFieldMap, err := helper.GetSQLFieldMapFromEntity(&entity.User{})
				if err != nil {
					t.Fatalf("an error '%s' was not expected when getting sqlFieldMap", err)
				}
				mock.ExpectQuery("SELECT (.+) users").
					WillReturnRows(sqlmock.NewRows(sqlFieldMap.Keys).
						AddRow("student-01", "FirstName 01", "LastName 01", "student-01@example.com", time.Now(), time.Now(), nil).
						RowError(0, sqlmock.ErrCancelled),
					)
			},
		},
	}

	for _, testCase := range testCases {
		testCase.Setup(ctx)
		_, err := repo.List(ctx, testCase.Req.(*ListParams))
		if testCase.ExpectedErr != nil {
			assert.Equal(t, testCase.ExpectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.ExpectedErr, err)
		}
	}
}

func TestUserRepo_Upsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := UserRepo{DB: db}

	testCases := []helper.TestCase{
		{
			Name: "happy case",
			Req:  &entity.User{},
			Setup: func(ctx context.Context) {
				mock.ExpectExec("INSERT INTO users").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			Name:        "error: bad connection",
			Req:         &entity.User{},
			ExpectedErr: sql.ErrConnDone,
			Setup: func(ctx context.Context) {
				mock.ExpectExec("INSERT INTO users").
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, testCase := range testCases {
		testCase.Setup(ctx)
		err := repo.Upsert(ctx, testCase.Req.(*entity.User))
		if testCase.ExpectedErr != nil {
			assert.Equal(t, testCase.ExpectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.ExpectedErr, err)
		}
	}
}
