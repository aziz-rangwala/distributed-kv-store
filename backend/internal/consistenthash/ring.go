package consistenthash

import (
	"hash/crc32"
	"sort"
	"fmt"
)

type Ring struct {
	hashMap  map[uint32]string  // Maps hash values to node addresses
	keys     []uint32           // Sorted list of hash values
	replicas int                // Number of virtual nodes per physical node
}

func NewRing(replicas int) *Ring {
	return &Ring{
		hashMap:  make(map[uint32]string),
		replicas: replicas,
	}
}

// AddNode adds a node to the ring with virtual replicas
func (r *Ring) AddNode(node string) {
	for i := 0; i < r.replicas; i++ {
		virtualNode := fmt.Sprintf("%s-%d", node, i)
		hash := hashKey(virtualNode)
		r.keys = append(r.keys, hash)
		r.hashMap[hash] = node
	}
	sort.Slice(r.keys, func(i, j int) bool { return r.keys[i] < r.keys[j] })
}

// GetNode returns the node responsible for a key
func (r *Ring) GetNode(key string) string {
	if len(r.keys) == 0 {
		return ""
	}
	hash := hashKey(key)
	idx := sort.Search(len(r.keys), func(i int) bool { return r.keys[i] >= hash })
	if idx == len(r.keys) {
		idx = 0
	}
	return r.hashMap[r.keys[idx]]
}

// GetReplicasCount returns the number of replicas per node
func (r *Ring) GetReplicasCount() int {
	return r.replicas
}

// GetNextNode returns the next node in the ring after the specified node
func (r *Ring) GetNextNode(node string) string {
	if len(r.keys) == 0 {
		return ""
	}

	// Find the current node's position
	var currentIdx int
	for idx, hash := range r.keys {
		if r.hashMap[hash] == node {
			currentIdx = idx
			break
		}
	}

	// Get the next node (wrap around if needed)
	nextIdx := (currentIdx + 1) % len(r.keys)
	return r.hashMap[r.keys[nextIdx]]
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}