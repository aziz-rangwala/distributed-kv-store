# Distributed Key-Value Store  

A fault-tolerant key-value storage system using **Raft consensus**

## Run with Docker  
```bash
git clone https://github.com/yourusername/distributed-kv-store.git  
cd distributed-kv-store
cd deployments  
docker-compose up --build
```

Write:

```bash
curl -X POST http://localhost:8081/keys \
  -H "Content-Type: application/json" \
  -d '{"key": "city", "value": "tokyo"}'
```

Read (any node):

```bash
curl http://localhost:8080/keys/city
```  
