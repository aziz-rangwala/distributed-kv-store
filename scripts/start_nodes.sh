#!/bin/bash

# Start 3 backend nodes
for port in {8080..8082}; do
  docker run -d \
    -p $port:9090 \
    -e NODE_ID=$((port - 8079)) \
    -e RAFT_ADDRESS="node$((port - 8079)):9090" \
    distributed-kv-store
done

echo "Nodes started on ports 8080-8082"