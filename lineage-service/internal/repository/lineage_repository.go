package repository

import (
	"context"
	"database/sql"
)

type LineageEvent struct {
	ID          uint64
	SubjectID   uint64
	Operation   string
	Source      string
	Destination string
	Purpose     string
	ConsentID   *uint64
	PayloadJSON string
	CreatedAt   string
}

type LineageRepository interface {
	Record(ctx context.Context, e LineageEvent) (uint64, error)
	ListBySubject(ctx context.Context, subjectID uint64) ([]LineageEvent, error)
}

type mysqlLineageRepository struct {
	db *sql.DB
}

func NewLineageRepository(db *sql.DB) LineageRepository {
	return &mysqlLineageRepository{db: db}
}

func (r *mysqlLineageRepository) Record(ctx context.Context, e LineageEvent) (uint64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO data_lineage (subject_id, operation, source, destination, purpose, consent_id, payload_json)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, e.SubjectID, e.Operation, e.Source, e.Destination, e.Purpose, e.ConsentID, e.PayloadJSON)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (r *mysqlLineageRepository) ListBySubject(ctx context.Context, subjectID uint64) ([]LineageEvent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, subject_id, operation, source, destination, purpose, consent_id, payload_json, created_at
		FROM data_lineage
		WHERE subject_id = ?
		ORDER BY created_at
	`, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]LineageEvent, 0)
	for rows.Next() {
		var e LineageEvent
		if err := rows.Scan(&e.ID, &e.SubjectID, &e.Operation, &e.Source, &e.Destination, &e.Purpose, &e.ConsentID, &e.PayloadJSON, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
