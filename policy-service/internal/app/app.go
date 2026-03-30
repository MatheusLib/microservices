package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"policy-service/internal/handler"
	"policy-service/internal/repository"
	"policy-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewPolicyRepository(db)
	svc := service.NewPolicyService(repo)
	h := handler.NewPolicyHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", handler.Health)
	r.Get("/policies", h.List)
	r.Post("/policies", h.Create)
	return r
}
