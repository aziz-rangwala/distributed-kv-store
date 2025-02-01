package raft

import (
	"encoding/json"
	"io"
	"sync"
	"fmt"
	"github.com/hashicorp/raft"
)

type RaftState struct {
	data map[string]string
	mu   sync.RWMutex
}

func (s *RaftState) Apply(l *raft.Log) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	var cmd map[string]string
	if err := json.Unmarshal(l.Data, &cmd); err != nil {
		return err
	}

	switch cmd["op"] {
	case "set":
		s.data[cmd["key"]] = cmd["value"]
	case "delete":
		delete(s.data, cmd["key"])
	default:
		return fmt.Errorf("unknown operation: %s", cmd["op"])
	}
	return nil
}

func (s *RaftState) Snapshot() (raft.FSMSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make(map[string]string)
	for k, v := range s.data {
		snapshot[k] = v
	}

	return &Snapshot{data: snapshot}, nil
}

func (s *RaftState) Restore(rc io.ReadCloser) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer rc.Close()

	s.data = make(map[string]string)
	return json.NewDecoder(rc).Decode(&s.data)
}

type Snapshot struct {
	data map[string]string
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	err := json.NewEncoder(sink).Encode(s.data)
	if err != nil {
		sink.Cancel()
		return err
	}
	return sink.Close()
}

func (s *Snapshot) Release() {}