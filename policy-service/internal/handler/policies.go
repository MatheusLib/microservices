package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"policy-service/internal/repository"
	"policy-service/internal/service"
)

type Policy struct {
	ID          uint64 `json:"id"`
	Version     string `json:"version"`
	ContentHash string `json:"content_hash"`
}

type createPolicyRequest struct {
	Version     string `json:"version"`
	ContentHash string `json:"content_hash"`
}

type PolicyHandler struct {
	Service service.PolicyService
}

func NewPolicyHandler(svc service.PolicyService) PolicyHandler {
	return PolicyHandler{Service: svc}
}

func (h PolicyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 1 || parsed > 1000 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	policies, err := h.Service.ListPolicies(ctx, limit)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}

	resp := make([]Policy, 0, len(policies))
	for _, p := range policies {
		resp = append(resp, Policy{
			ID:          p.ID,
			Version:     p.Version,
			ContentHash: p.ContentHash,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h PolicyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req createPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.Version == "" || req.ContentHash == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreatePolicy(ctx, repository.Policy{
		Version:     req.Version,
		ContentHash: req.ContentHash,
	})
	if err != nil {
		http.Error(w, "create error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]uint64{"id": id})
}
