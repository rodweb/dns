package api

import (
	"net/http"
)

type RecordService interface {
	GetAll(recordType string, name string) ([]*interface{}, error)
	Create(recordType string, name string, value string, ttl int, note string) error
	Delete(recordType string, name string) error
}

type Handler struct {
	recordService RecordService
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleRecords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listRecords(w, r)
	case http.MethodPost:
		h.createRecord(w, r)
	case http.MethodPut:
		h.updateRecord(w, r)
	case http.MethodDelete:
		h.deleteRecord(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listRecords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) updateRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) deleteRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func parsePath(path string) (recordType string, name string) {
	return
}
