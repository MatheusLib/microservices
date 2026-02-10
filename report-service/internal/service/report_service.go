package service

import (
	"context"

	"report-service/internal/repository"
)

type ReportService interface {
	ListConsents(ctx context.Context, userID *uint64, limit int) ([]repository.ConsentReport, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) ListConsents(ctx context.Context, userID *uint64, limit int) ([]repository.ConsentReport, error) {
	return s.repo.ListConsents(ctx, userID, limit)
}
