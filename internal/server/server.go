package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"AP_WebServer/internal/model"
	"AP_WebServer/internal/store"
)

type Server struct {
	store     *store.Store[string, string]
	startTime time.Time
	requests  int64
}

func NewServer(store *store.Store[string, string]) *Server {
	return &Server{
		store:     store,
		startTime: time.Now(),
	}
}

func (s *Server) RequestCount() int64 {
	return atomic.LoadInt64(&s.requests)
}

func (s *Server) KeyCount() int {
	return s.store.Count()
}

func (s *Server) incrementRequests() {
	atomic.AddInt64(&s.requests, 1)
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/data", s.handleData)
	mux.HandleFunc("/data/", s.handleDataByKey)
	mux.HandleFunc("/stats", s.handleStats)
	return mux
}

// handlers
func (s *Server) handleData(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()

	switch r.Method {
	case http.MethodPost:
		var body struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil ||
			body.Key == "" || body.Value == "" {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		s.store.Set(body.Key, body.Value)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)

	case http.MethodGet:
		json.NewEncoder(w).Encode(s.store.Snapshot())

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleDataByKey(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()
	key := strings.TrimPrefix(r.URL.Path, "/data/")

	switch r.Method {
	case http.MethodGet:
		if val, ok := s.store.Get(key); ok {
			json.NewEncoder(w).Encode(map[string]string{key: val})
			return
		}
		w.WriteHeader(http.StatusNotFound)

	case http.MethodDelete:
		if s.store.Delete(key) {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()

	stats := model.Stats{
		Requests:      s.RequestCount(),
		Keys:          s.KeyCount(),
		UptimeSeconds: int64(time.Since(s.startTime).Seconds()),
	}

	json.NewEncoder(w).Encode(stats)
}
