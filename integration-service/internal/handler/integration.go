package handler

import (
	"context"
	"net/http"
	"time"

	"integration-service/internal/external_services"
)

type IntegrationHandler struct {
	Service external_services.Service
}

func NewIntegrationHandler(svc external_services.Service) IntegrationHandler {
	return IntegrationHandler{Service: svc}
}

func (h IntegrationHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.Service.Ping(ctx); err != nil {
		http.Error(w, "external unreachable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
