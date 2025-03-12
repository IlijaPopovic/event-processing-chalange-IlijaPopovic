package metrics

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	storage *Storage
}

func NewServer(storage *Storage) *Server {
	return &Server{storage: storage}
}

func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/materialized", s.handleGetMetrics).Methods("GET")
}

func (s *Server) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metrics := s.storage.GetMetrics()
	json.NewEncoder(w).Encode(metrics)
}