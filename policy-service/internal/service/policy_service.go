package service

import (
	"context"

	"policy-service/internal/repository"
)

type PolicyService interface {
	ListPolicies(ctx context.Context, limit int) ([]repository.Policy, error)
	CreatePolicy(ctx context.Context, p repository.Policy) (uint64, error)
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

func (s *policyService) CreatePolicy(ctx context.Context, p repository.Policy) (uint64, error) {
	return s.repo.Create(ctx, p)
}
