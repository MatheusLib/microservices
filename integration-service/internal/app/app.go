package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"integration-service/internal/external_services"
	"integration-service/internal/handler"
)

func NewRouter(baseURL string) http.Handler {
	client := external_services.NewHTTPClient(baseURL)
	svc := external_services.NewService(client)
	h := handler.NewIntegrationHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", handler.Health)
	r.Get("/integrations/ping", h.Ping)
	return r
}
