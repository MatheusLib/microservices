package app

import (
	"database/sql"
	"net/http"

	"consent-service/internal/handler"
	"consent-service/internal/repository"
	"consent-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewConsentRepository(db)
	svc := service.NewConsentService(repo)
	h := handler.NewConsentHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/consents", h.List)
	return mux
}
