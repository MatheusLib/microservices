package repository

import (
	"context"
	"database/sql"
)

type ConsentReport struct {
	ID       uint64
	UserID   uint64
	PolicyID uint64
	Purpose  string
	Status   string
}

type ReportRepository interface {
	ListConsents(ctx context.Context, userID *uint64, limit int) ([]ConsentReport, error)
}

type mysqlReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &mysqlReportRepository{db: db}
}

func (r *mysqlReportRepository) ListConsents(ctx context.Context, userID *uint64, limit int) ([]ConsentReport, error) {
	if userID == nil {
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

		results := make([]ConsentReport, 0)
		for rows.Next() {
			var c ConsentReport
			if err := rows.Scan(&c.ID, &c.UserID, &c.PolicyID, &c.Purpose, &c.Status); err != nil {
				return nil, err
			}
			results = append(results, c)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return results, nil
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, policy_id, purpose, status
		FROM consents
		WHERE user_id = ?
		ORDER BY id
		LIMIT ?
	`, *userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]ConsentReport, 0)
	for rows.Next() {
		var c ConsentReport
		if err := rows.Scan(&c.ID, &c.UserID, &c.PolicyID, &c.Purpose, &c.Status); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
