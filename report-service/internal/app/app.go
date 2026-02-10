package app

import (
	"database/sql"
	"net/http"

	"report-service/internal/handler"
	"report-service/internal/repository"
	"report-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewReportRepository(db)
	svc := service.NewReportService(repo)
	h := handler.NewReportHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/reports/consents", h.ListConsents)
	return mux
}
