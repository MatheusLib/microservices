package service

import (
	"context"

	"lineage-service/internal/repository"
)

type LineageService interface {
	Record(ctx context.Context, e repository.LineageEvent) (uint64, error)
	ExportBySubject(ctx context.Context, subjectID uint64) ([]repository.LineageEvent, error)
}

type lineageService struct {
	repo repository.LineageRepository
}

func NewLineageService(repo repository.LineageRepository) LineageService {
	return &lineageService{repo: repo}
}

func (s *lineageService) Record(ctx context.Context, e repository.LineageEvent) (uint64, error) {
	return s.repo.Record(ctx, e)
}

func (s *lineageService) ExportBySubject(ctx context.Context, subjectID uint64) ([]repository.LineageEvent, error) {
	return s.repo.ListBySubject(ctx, subjectID)
}
