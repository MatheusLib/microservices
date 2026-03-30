package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestConfigure_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := configure(ctx, db)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil db")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestConfigure_PingError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(context.DeadlineExceeded)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = configure(ctx, db)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewMySQL_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	original := openFunc
	openFunc = func(driver, dsn string) (*sql.DB, error) {
		return db, nil
	}
	defer func() { openFunc = original }()

	mock.ExpectPing()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := NewMySQL(ctx, "localhost", "3306", "user", "pass", "testdb")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil db")
	}
}

func TestNewMySQL_OpenError(t *testing.T) {
	original := openFunc
	openFunc = func(driver, dsn string) (*sql.DB, error) {
		return nil, errors.New("open error")
	}
	defer func() { openFunc = original }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := NewMySQL(ctx, "localhost", "3306", "user", "pass", "testdb")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewMySQL_PingError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	original := openFunc
	openFunc = func(driver, dsn string) (*sql.DB, error) {
		return db, nil
	}
	defer func() { openFunc = original }()

	mock.ExpectPing().WillReturnError(errors.New("ping fail"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = NewMySQL(ctx, "localhost", "3306", "user", "pass", "testdb")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
