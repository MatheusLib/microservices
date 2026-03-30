package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestConsentList_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "user_id", "policy_id", "purpose", "status"}).
		AddRow(1, 10, 2, "marketing", "active")
	mock.ExpectQuery("SELECT id, user_id").WillReturnRows(rows)
	repo := NewConsentRepository(db)
	result, err := repo.List(context.Background(), 10)
	if err != nil || len(result) != 1 {
		t.Fatalf("unexpected: err=%v len=%d", err, len(result))
	}
}

func TestConsentList_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT id, user_id").WillReturnError(errors.New("db down"))
	repo := NewConsentRepository(db)
	_, err := repo.List(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConsentList_RowsErr(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "user_id", "policy_id", "purpose", "status"}).
		AddRow(1, 10, 2, "x", "active").
		RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT id, user_id").WillReturnRows(rows)
	repo := NewConsentRepository(db)
	_, err := repo.List(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConsentCreate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO consents").
		WithArgs(uint64(1), uint64(2), "marketing", "active").
		WillReturnResult(sqlmock.NewResult(5, 1))
	repo := NewConsentRepository(db)
	id, err := repo.Create(context.Background(), Consent{UserID: 1, PolicyID: 2, Purpose: "marketing", Status: "active"})
	if err != nil || id != 5 {
		t.Fatalf("unexpected: err=%v id=%d", err, id)
	}
}

func TestConsentCreate_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO consents").WillReturnError(errors.New("fail"))
	repo := NewConsentRepository(db)
	_, err := repo.Create(context.Background(), Consent{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConsentRevoke_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("UPDATE consents").
		WithArgs("3").
		WillReturnResult(sqlmock.NewResult(0, 1))
	repo := NewConsentRepository(db)
	if err := repo.Revoke(context.Background(), "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConsentRevoke_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("UPDATE consents").WillReturnError(errors.New("fail"))
	repo := NewConsentRepository(db)
	if err := repo.Revoke(context.Background(), "1"); err == nil {
		t.Fatal("expected error")
	}
}
