package app

import (
	"database/sql"
	"net/http"

	"audit-service/internal/handler"
	"audit-service/internal/repository"
	"audit-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewAuditRepository(db)
	svc := service.NewAuditService(repo)
	h := handler.NewAuditHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/audit/events", h.List)
	return mux
}
