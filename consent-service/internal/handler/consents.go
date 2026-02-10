package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"consent-service/internal/service"
)

type Consent struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"user_id"`
	PolicyID uint64 `json:"policy_id"`
	Purpose  string `json:"purpose"`
	Status   string `json:"status"`
}

type ConsentHandler struct {
	Service service.ConsentService
}

func NewConsentHandler(svc service.ConsentService) ConsentHandler {
	return ConsentHandler{Service: svc}
}

func (h ConsentHandler) List(w http.ResponseWriter, r *http.Request) {
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

	consents, err := h.Service.ListConsents(ctx, limit)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}

	resp := make([]Consent, 0, len(consents))
	for _, c := range consents {
		resp = append(resp, Consent{
			ID:       c.ID,
			UserID:   c.UserID,
			PolicyID: c.PolicyID,
			Purpose:  c.Purpose,
			Status:   c.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
