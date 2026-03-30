package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"report-service/internal/handler"
	"report-service/internal/repository"
	"report-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewReportRepository(db)
	svc := service.NewReportService(repo)
	h := handler.NewReportHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", handler.Health)
	r.Get("/reports/consents", h.ListConsents)
	return r
}
