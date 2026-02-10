package app

import (
	"database/sql"
	"net/http"

	"policy-service/internal/handler"
	"policy-service/internal/repository"
	"policy-service/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	repo := repository.NewPolicyRepository(db)
	svc := service.NewPolicyService(repo)
	h := handler.NewPolicyHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/policies", h.List)
	return mux
}
