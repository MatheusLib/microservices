package service

import (
	"context"

	"policy-service/internal/repository"
)

type PolicyService interface {
	ListPolicies(ctx context.Context, limit int) ([]repository.Policy, error)
}

type policyService struct {
	repo repository.PolicyRepository
}

func NewPolicyService(repo repository.PolicyRepository) PolicyService {
	return &policyService{repo: repo}
}

func (s *policyService) ListPolicies(ctx context.Context, limit int) ([]repository.Policy, error) {
	return s.repo.List(ctx, limit)
}
