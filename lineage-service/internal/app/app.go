package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"lineage-service/internal/handler"
	"lineage-service/internal/repository"
	"lineage-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewLineageRepository(db)
	svc := service.NewLineageService(repo)
	h := handler.NewLineageHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", handler.Health)
	r.Post("/lineage", h.Record)
	r.Get("/lineage/export/{subject_id}", h.Export)
	return r
}
