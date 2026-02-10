package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"report-service/internal/service"
)

type ConsentReport struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"user_id"`
	PolicyID uint64 `json:"policy_id"`
	Purpose  string `json:"purpose"`
	Status   string `json:"status"`
}

type ReportHandler struct {
	Service service.ReportService
}

func NewReportHandler(svc service.ReportService) ReportHandler {
	return ReportHandler{Service: svc}
}

func (h ReportHandler) ListConsents(w http.ResponseWriter, r *http.Request) {
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

	var userID *uint64
	if v := r.URL.Query().Get("user_id"); v != "" {
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		userID = &parsed
	}

	reports, err := h.Service.ListConsents(ctx, userID, limit)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}

	resp := make([]ConsentReport, 0, len(reports))
	for _, c := range reports {
		resp = append(resp, ConsentReport{
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
