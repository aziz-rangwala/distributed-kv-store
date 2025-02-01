#!/bin/bash

# Test set operations
for i in {1..100}; do
  curl -X POST http://localhost:8080/keys \
    -H "Content-Type: application/json" \
    -d "{\"key\": \"key$i\", \"value\": \"value$i\"}"
done

# Test get operations
for i in {1..100}; do
  curl http://localhost:8080/keys/key$i
  sleep 0.1
done

# Use k6 for advanced load testing
k6 run --vus 10 --duration 30s <(echo '
import http from "k6/http";
export default function() {
  http.post("http://localhost:8080/keys", JSON.stringify({
    key: `key${__VU}`,
    value: "value"
  }));
}
')