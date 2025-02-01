package store

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hashicorp/raft"
)

const (
	applyTimeout = 10 * time.Second
)

type Store struct {
	mu         sync.RWMutex
	data       map[string]string
	Raft       *raft.Raft // Directly use HashiCorp's Raft
	localAddr  string
}

func NewStore(raftNode *raft.Raft, localAddr string) *Store {
	return &Store{
		data:      make(map[string]string),
		Raft:      raftNode,
		localAddr: localAddr,
	}
}

// Get retrieves a value from the store
func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}
	return value, nil
}

// Set stores a value using Raft consensus
func (s *Store) Set(key, value string) error {
	if s.Raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	cmd := map[string]string{"op": "set", "key": key, "value": value}
	b, _ := json.Marshal(cmd)
	return s.Raft.Apply(b, applyTimeout).Error()
}

// Implement raft.FSM interface
func (s *Store) Apply(l *raft.Log) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	var cmd map[string]string
	json.Unmarshal(l.Data, &cmd)
	switch cmd["op"] {
	case "set":
		s.data[cmd["key"]] = cmd["value"]
	case "delete":
		delete(s.data, cmd["key"])
	}
	return nil
}

func (s *Store) Snapshot() (raft.FSMSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot := make(map[string]string)
	for k, v := range s.data {
		snapshot[k] = v
	}
	return &Snapshot{data: snapshot}, nil
}

func (s *Store) Restore(rc io.ReadCloser) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer rc.Close()
	return json.NewDecoder(rc).Decode(&s.data)
}

type Snapshot struct{ data map[string]string }
func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	return json.NewEncoder(sink).Encode(s.data)
}
func (s *Snapshot) Release() {}