package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"consent-service/internal/handler"
	"consent-service/internal/repository"
	"consent-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewConsentRepository(db)
	svc := service.NewConsentService(repo)
	h := handler.NewConsentHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", handler.Health)
	r.Get("/consents", h.List)
	r.Post("/consents", h.Create)
	r.Patch("/consents/{document_id}/revoke", h.Revoke)
	return r
}
