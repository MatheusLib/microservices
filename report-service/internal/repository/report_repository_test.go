package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func u64(v uint64) *uint64 { return &v }

func TestReportListConsents_NoFilter(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "user_id", "policy_id", "purpose", "status"}).
		AddRow(1, 10, 2, "marketing", "active")
	mock.ExpectQuery("SELECT id, user_id").WillReturnRows(rows)
	repo := NewReportRepository(db)
	result, err := repo.ListConsents(context.Background(), nil, 10)
	if err != nil || len(result) != 1 {
		t.Fatalf("unexpected: err=%v len=%d", err, len(result))
	}
}

func TestReportListConsents_WithFilter(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "user_id", "policy_id", "purpose", "status"}).
		AddRow(1, 10, 2, "marketing", "active")
	mock.ExpectQuery("SELECT id, user_id").WillReturnRows(rows)
	repo := NewReportRepository(db)
	result, err := repo.ListConsents(context.Background(), u64(10), 10)
	if err != nil || len(result) != 1 {
		t.Fatalf("unexpected: err=%v len=%d", err, len(result))
	}
}

func TestReportListConsents_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT id, user_id").WillReturnError(errors.New("db down"))
	repo := NewReportRepository(db)
	_, err := repo.ListConsents(context.Background(), nil, 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReportListConsents_DBErrorWithFilter(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT id, user_id").WillReturnError(errors.New("db down"))
	repo := NewReportRepository(db)
	_, err := repo.ListConsents(context.Background(), u64(5), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}
