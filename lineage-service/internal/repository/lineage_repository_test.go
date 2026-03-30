package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRecord_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	consentID := uint64(42)
	mock.ExpectExec("INSERT INTO data_lineage").
		WithArgs(uint64(1), "COLLECT", "api", "storage", "analytics", &consentID, "{}").
		WillReturnResult(sqlmock.NewResult(5, 1))

	repo := NewLineageRepository(db)
	id, err := repo.Record(context.Background(), LineageEvent{
		SubjectID:   1,
		Operation:   "COLLECT",
		Source:      "api",
		Destination: "storage",
		Purpose:     "analytics",
		ConsentID:   &consentID,
		PayloadJSON: "{}",
	})
	if err != nil || id != 5 {
		t.Fatalf("unexpected: err=%v id=%d", err, id)
	}
}

func TestRecord_NilConsentID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("INSERT INTO data_lineage").
		WithArgs(uint64(1), "DELETE", "api", "storage", "compliance", nil, "{}").
		WillReturnResult(sqlmock.NewResult(7, 1))

	repo := NewLineageRepository(db)
	id, err := repo.Record(context.Background(), LineageEvent{
		SubjectID:   1,
		Operation:   "DELETE",
		Source:      "api",
		Destination: "storage",
		Purpose:     "compliance",
		ConsentID:   nil,
		PayloadJSON: "{}",
	})
	if err != nil || id != 7 {
		t.Fatalf("unexpected: err=%v id=%d", err, id)
	}
}

func TestRecord_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO data_lineage").WillReturnError(errors.New("fail"))
	repo := NewLineageRepository(db)
	_, err := repo.Record(context.Background(), LineageEvent{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListBySubject_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	consentID := uint64(10)
	rows := sqlmock.NewRows([]string{"id", "subject_id", "operation", "source", "destination", "purpose", "consent_id", "payload_json", "created_at"}).
		AddRow(1, 100, "COLLECT", "api", "db", "analytics", &consentID, "{}", "2025-01-01 00:00:00").
		AddRow(2, 100, "SHARE", "db", "partner", "marketing", nil, "{}", "2025-01-02 00:00:00")
	mock.ExpectQuery("SELECT id, subject_id").WillReturnRows(rows)

	repo := NewLineageRepository(db)
	result, err := repo.ListBySubject(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
	if result[0].Operation != "COLLECT" {
		t.Errorf("expected COLLECT, got %s", result[0].Operation)
	}
	if result[1].ConsentID != nil {
		t.Errorf("expected nil consent_id for second event")
	}
}

func TestListBySubject_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT id, subject_id").WillReturnError(errors.New("db down"))
	repo := NewLineageRepository(db)
	_, err := repo.ListBySubject(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListBySubject_RowsErr(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "subject_id", "operation", "source", "destination", "purpose", "consent_id", "payload_json", "created_at"}).
		AddRow(1, 100, "COLLECT", "api", "db", "analytics", nil, "{}", "2025-01-01").
		RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT id, subject_id").WillReturnRows(rows)
	repo := NewLineageRepository(db)
	_, err := repo.ListBySubject(context.Background(), 100)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListBySubject_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "subject_id", "operation", "source", "destination", "purpose", "consent_id", "payload_json", "created_at"})
	mock.ExpectQuery("SELECT id, subject_id").WillReturnRows(rows)
	repo := NewLineageRepository(db)
	result, err := repo.ListBySubject(context.Background(), 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 events, got %d", len(result))
	}
}
