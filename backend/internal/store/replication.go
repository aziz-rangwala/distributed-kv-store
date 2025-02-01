package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"distributed-kv-store/backend/internal/consistenthash"
)

const (
	replicationTimeout = 2 * time.Second
)

type Replicator struct {
	ring       *consistenthash.Ring
	localAddr  string
	httpClient *http.Client
	mu         sync.Mutex
}

func NewReplicator(ring *consistenthash.Ring, localAddr string) *Replicator {
	return &Replicator{
		ring:       ring,
		localAddr:  localAddr,
		httpClient: &http.Client{Timeout: replicationTimeout},
	}
}

// Replicate sends a key-value pair to replica nodes
func (r *Replicator) Replicate(key, value string) error {
	replicas := r.GetReplicas(key)
	if len(replicas) == 0 {
		return fmt.Errorf("no replicas found for key: %s", key)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(replicas))

	for _, replica := range replicas {
		if replica == r.localAddr {
			continue // Skip self
		}

		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			if err := r.sendToReplica(addr, key, value); err != nil {
				log.Printf("Replication to %s failed: %v", addr, err)
				errCh <- fmt.Errorf("replica %s: %v", addr, err)
			}
		}(replica)
	}

	wg.Wait()
	close(errCh)

	// Collect errors
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("replication errors: %v", errs)
	}

	return nil
}

// sendToReplica sends a key-value pair to a specific replica node
func (r *Replicator) sendToReplica(addr, key, value string) error {
	url := fmt.Sprintf("http://%s/keys", addr)
	payload := map[string]string{
		"key":   key,
		"value": value,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := r.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetReplicas returns the list of replica nodes for a key
func (r *Replicator) GetReplicas(key string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	replicas := make([]string, 0)
	primaryNode := r.ring.GetNode(key)

	// Add primary node
	replicas = append(replicas, primaryNode)

	// Add additional replicas (next nodes in the ring)
	for i := 1; i <= r.ring.GetReplicasCount(); i++ {
		nextNode := r.ring.GetNextNode(primaryNode)
		if nextNode == "" {
			break
		}
		replicas = append(replicas, nextNode)
		primaryNode = nextNode
	}

	return replicas
}