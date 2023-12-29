package api

import "net/http"

func registerRoutes() {
	h := NewHandler()
	http.HandleFunc("/records", h.HandleRecords)
}
