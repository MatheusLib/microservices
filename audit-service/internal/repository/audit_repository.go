package repository

import (
	"context"
	"database/sql"
)

type AuditEvent struct {
	ID         uint64
	EventType  string
	EntityType string
	EntityID   uint64
	Payload    string
}

type AuditRepository interface {
	List(ctx context.Context, limit int) ([]AuditEvent, error)
}

type mysqlAuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) AuditRepository {
	return &mysqlAuditRepository{db: db}
}

func (r *mysqlAuditRepository) List(ctx context.Context, limit int) ([]AuditEvent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, event_type, entity_type, entity_id, payload_json
		FROM audit_events
		ORDER BY id
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]AuditEvent, 0)
	for rows.Next() {
		var e AuditEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.EntityType, &e.EntityID, &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
