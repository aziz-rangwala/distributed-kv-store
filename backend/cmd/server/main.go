package main

import (
    "net/http"
    "log"
    "net"
    "os"
    "path/filepath"
    "time"
    "strings"
    "distributed-kv-store/backend/internal/api"
    internalStore "distributed-kv-store/backend/internal/store"
    "github.com/hashicorp/raft"
    raftboltdb "github.com/hashicorp/raft-boltdb"
)

func main() {
    // Initialize Raft configuration
    config := raft.DefaultConfig()
    nodeID := os.Getenv("NODE_ID")
    config.LocalID = raft.ServerID(nodeID)

    // Get Raft address and port
    raftAddr := os.Getenv("RAFT_ADDRESS")
    addr, err := net.ResolveTCPAddr("tcp", raftAddr)
    if err != nil {
        log.Fatal("Failed to resolve TCP address:", err)
    }

    // Create Raft transport
    transport, err := raft.NewTCPTransport(raftAddr, addr, 3, 10*time.Second, os.Stderr)
    if err != nil {
        log.Fatal("Failed to create transport:", err)
    }

    // Raft storage directory
    raftDir := filepath.Join("raft", nodeID)
    if err := os.MkdirAll(raftDir, 0700); err != nil {
        log.Fatal("Failed to create raft directory:", err)
    }

    // BoltDB store for Raft logs
    raftStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft.db"))
    if err != nil {
        log.Fatal("Failed to create BoltStore:", err)
    }

    // Snapshot store
    snapshots, err := raft.NewFileSnapshotStore(raftDir, 3, os.Stderr)
    if err != nil {
        log.Fatal("Failed to create snapshot store:", err)
    }

    // State machine
    fsm := internalStore.NewStore(nil, "")

    // Create Raft node
    raftNode, err := raft.NewRaft(config, fsm, raftStore, raftStore, snapshots, transport)
    if err != nil {
        log.Fatal("Raft init failed:", err)
    }

    // Check if the node has existing state (like logs or snapshots)
    hasState, err := raft.HasExistingState(raftStore, raftStore, snapshots)
    if err != nil {
        log.Fatal("Error checking for existing state:", err)
    }

    if !hasState {
        // Parse peer list from environment
        peers := os.Getenv("RAFT_PEERS")
        servers := []raft.Server{}
        for _, peer := range strings.Split(peers, ",") {
            // Extract node ID from "nodeX:9090"
            nodeIDFromPeer := strings.Split(peer, ":")[0]
            servers = append(servers, raft.Server{
                ID:      raft.ServerID(nodeIDFromPeer),
                Address: raft.ServerAddress(peer),
            })
        }

        // Bootstrap the cluster with all peers
        configFuture := raftNode.BootstrapCluster(raft.Configuration{
            Servers: servers,
        })
        if err := configFuture.Error(); err != nil {
            log.Printf("Bootstrap failed (normal if not the first node): %v", err)
        }
    }


    // Start HTTP server
    httpPort := os.Getenv("HTTP_PORT")
    router := api.NewRouter(fsm, raftNode)
    log.Fatal(http.ListenAndServe(":"+httpPort, router))
}