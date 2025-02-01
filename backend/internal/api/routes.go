package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"log"
	"github.com/hashicorp/raft"
	"distributed-kv-store/backend/internal/store"
)

type Server struct {
	store *store.Store
	raft  *raft.Raft // Directly use HashiCorp's Raft
}

func NewRouter(store *store.Store, raftNode *raft.Raft) *http.ServeMux {
	s := &Server{store: store, raft: raftNode}
	mux := http.NewServeMux()
	mux.HandleFunc("/keys/", s.handleGet)
	mux.HandleFunc("/keys", s.handleSet)
	return mux
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/keys/")
	value, err := s.store.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"value": value})
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
    var req struct{ Key, Value string }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if s.raft.State() != raft.Leader {
        leaderAddr := string(s.raft.Leader())
        http.Redirect(w, r, fmt.Sprintf("http://%s/keys", leaderAddr), http.StatusTemporaryRedirect)
        return
    }

    if s.store == nil {
        log.Println("Store is nil")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if err := s.store.Set(req.Key, req.Value); err != nil {
        log.Printf("Failed to set key-value pair: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}