package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPolicyList_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "version", "content_hash"}).
		AddRow(1, "v1", "abc123")
	mock.ExpectQuery("SELECT id, version").WillReturnRows(rows)
	repo := NewPolicyRepository(db)
	result, err := repo.List(context.Background(), 10)
	if err != nil || len(result) != 1 {
		t.Fatalf("unexpected: err=%v len=%d", err, len(result))
	}
}

func TestPolicyList_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT id, version").WillReturnError(errors.New("db down"))
	repo := NewPolicyRepository(db)
	_, err := repo.List(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPolicyList_RowsErr(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "version", "content_hash"}).
		AddRow(1, "v1", "abc").
		RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT id, version").WillReturnRows(rows)
	repo := NewPolicyRepository(db)
	_, err := repo.List(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPolicyCreate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO policies").
		WithArgs("v2", "hash999").
		WillReturnResult(sqlmock.NewResult(7, 1))
	repo := NewPolicyRepository(db)
	id, err := repo.Create(context.Background(), Policy{Version: "v2", ContentHash: "hash999"})
	if err != nil || id != 7 {
		t.Fatalf("unexpected: err=%v id=%d", err, id)
	}
}

func TestPolicyCreate_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO policies").WillReturnError(errors.New("fail"))
	repo := NewPolicyRepository(db)
	_, err := repo.Create(context.Background(), Policy{})
	if err == nil {
		t.Fatal("expected error")
	}
}
