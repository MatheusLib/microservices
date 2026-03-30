package service

import (
	"context"

	"audit-service/internal/repository"
)

type AuditService interface {
	ListEvents(ctx context.Context, limit int) ([]repository.AuditEvent, error)
	RecordEvent(ctx context.Context, e repository.AuditEvent) (uint64, error)
}

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) ListEvents(ctx context.Context, limit int) ([]repository.AuditEvent, error) {
	return s.repo.List(ctx, limit)
}

func (s *auditService) RecordEvent(ctx context.Context, e repository.AuditEvent) (uint64, error) {
	return s.repo.Record(ctx, e)
}
