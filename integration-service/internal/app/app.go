package app

import (
	"net/http"

	"integration-service/internal/external_services"
	"integration-service/internal/handler"
)

func NewRouter(baseURL string) http.Handler {
	client := external_services.NewHTTPClient(baseURL)
	svc := external_services.NewService(client)
	h := handler.NewIntegrationHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/integrations/ping", h.Ping)
	return mux
}
