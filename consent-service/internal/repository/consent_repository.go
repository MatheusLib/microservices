package repository

import (
	"context"
	"database/sql"
)

type Consent struct {
	ID       uint64
	UserID   uint64
	PolicyID uint64
	Purpose  string
	Status   string
}

type ConsentRepository interface {
	List(ctx context.Context, limit int) ([]Consent, error)
}

type mysqlConsentRepository struct {
	db *sql.DB
}

func NewConsentRepository(db *sql.DB) ConsentRepository {
	return &mysqlConsentRepository{db: db}
}

func (r *mysqlConsentRepository) List(ctx context.Context, limit int) ([]Consent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, policy_id, purpose, status
		FROM consents
		ORDER BY id
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	consents := make([]Consent, 0)
	for rows.Next() {
		var c Consent
		if err := rows.Scan(&c.ID, &c.UserID, &c.PolicyID, &c.Purpose, &c.Status); err != nil {
			return nil, err
		}
		consents = append(consents, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return consents, nil
}
