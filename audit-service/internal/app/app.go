package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"audit-service/internal/handler"
	"audit-service/internal/repository"
	"audit-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewAuditRepository(db)
	svc := service.NewAuditService(repo)
	h := handler.NewAuditHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", handler.Health)
	r.Get("/audit/events", h.List)
	r.Post("/audit-events", h.Record)
	return r
}
