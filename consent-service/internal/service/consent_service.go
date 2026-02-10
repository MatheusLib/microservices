package service

import (
	"context"

	"consent-service/internal/repository"
)

type ConsentService interface {
	ListConsents(ctx context.Context, limit int) ([]repository.Consent, error)
}

type consentService struct {
	repo repository.ConsentRepository
}

func NewConsentService(repo repository.ConsentRepository) ConsentService {
	return &consentService{repo: repo}
}

func (s *consentService) ListConsents(ctx context.Context, limit int) ([]repository.Consent, error) {
	return s.repo.List(ctx, limit)
}
