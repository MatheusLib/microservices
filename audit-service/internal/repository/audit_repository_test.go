package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuditList_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "event_type", "entity_type", "entity_id", "payload_json"}).
		AddRow(1, "ConsentCreated", "consent", 1, `{"consent_id":1}`)
	mock.ExpectQuery("SELECT id, event_type").WillReturnRows(rows)
	repo := NewAuditRepository(db)
	result, err := repo.List(context.Background(), 10)
	if err != nil || len(result) != 1 {
		t.Fatalf("unexpected: err=%v len=%d", err, len(result))
	}
}

func TestAuditList_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT id, event_type").WillReturnError(errors.New("db down"))
	repo := NewAuditRepository(db)
	_, err := repo.List(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAuditList_RowsErr(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "event_type", "entity_type", "entity_id", "payload_json"}).
		AddRow(1, "ConsentCreated", "consent", 1, `{}`).
		RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT id, event_type").WillReturnRows(rows)
	repo := NewAuditRepository(db)
	_, err := repo.List(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAuditRecord_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO audit_events").
		WithArgs("ConsentCreated", "consent", uint64(1), `{}`).
		WillReturnResult(sqlmock.NewResult(3, 1))
	repo := NewAuditRepository(db)
	id, err := repo.Record(context.Background(), AuditEvent{
		EventType: "ConsentCreated", EntityType: "consent", EntityID: 1, Payload: `{}`,
	})
	if err != nil || id != 3 {
		t.Fatalf("unexpected: err=%v id=%d", err, id)
	}
}

func TestAuditRecord_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO audit_events").WillReturnError(errors.New("fail"))
	repo := NewAuditRepository(db)
	_, err := repo.Record(context.Background(), AuditEvent{})
	if err == nil {
		t.Fatal("expected error")
	}
}
