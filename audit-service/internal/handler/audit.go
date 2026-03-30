package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"audit-service/internal/repository"
	"audit-service/internal/service"
)

type AuditEvent struct {
	ID         uint64 `json:"id"`
	EventType  string `json:"event_type"`
	EntityType string `json:"entity_type"`
	EntityID   uint64 `json:"entity_id"`
	Payload    string `json:"payload_json"`
}

type recordRequest struct {
	EventType  string `json:"event_type"`
	EntityType string `json:"entity_type"`
	EntityID   uint64 `json:"entity_id"`
	Payload    string `json:"payload_json"`
}

type AuditHandler struct {
	Service service.AuditService
}

func NewAuditHandler(svc service.AuditService) AuditHandler {
	return AuditHandler{Service: svc}
}

func (h AuditHandler) List(w http.ResponseWriter, r *http.Request) {
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

	events, err := h.Service.ListEvents(ctx, limit)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}

	resp := make([]AuditEvent, 0, len(events))
	for _, e := range events {
		resp = append(resp, AuditEvent{
			ID:         e.ID,
			EventType:  e.EventType,
			EntityType: e.EntityType,
			EntityID:   e.EntityID,
			Payload:    e.Payload,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h AuditHandler) Record(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req recordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.EventType == "" || req.EntityType == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	id, err := h.Service.RecordEvent(ctx, repository.AuditEvent{
		EventType:  req.EventType,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		Payload:    req.Payload,
	})
	if err != nil {
		http.Error(w, "record error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]uint64{"id": id})
}
