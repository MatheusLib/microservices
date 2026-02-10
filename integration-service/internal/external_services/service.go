package external_services

import "context"

type Service interface {
	Ping(ctx context.Context) error
}

type service struct {
	client Client
}

func NewService(client Client) Service {
	return &service{client: client}
}

func (s *service) Ping(ctx context.Context) error {
	return s.client.Ping(ctx)
}
