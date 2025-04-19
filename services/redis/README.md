# 🔴 Redis

A simple Redis setup with persistence for development and testing purposes.

## 🚀 Usage

### Using InfraCLI (Recommended)

```bash
# Start the service
infracli run redis

# Get connection information
infracli info redis

# Stop the service
infracli down redis
```

### Manual Usage

```bash
cd redis
docker-compose up -d
```

## ⚙️ Configuration

- **Port:** 6379
- **Persistence:** Enabled (appendonly)
- **Connection String:** redis://localhost:6379

## 💾 Data Persistence

Data is stored in a named Docker volume `redis-data` which persists between container restarts. Redis is configured with append-only file (AOF) persistence for better durability.

To reset the database, use `infracli down redis --volumes` or `docker-compose down -v` to remove the volume as well.

## 🔧 Common Operations

### 💻 Connect using the Redis CLI:
```bash
docker exec -it redis redis-cli
```

### 🔍 Test connectivity:
```bash
docker exec -it redis redis-cli ping
```

### 📋 View all keys:
```bash
docker exec -it redis redis-cli keys "*"
```

### 📊 Get server info:
```bash
docker exec -it redis redis-cli info
```

### 🛡️ Set a key with expiration:
```bash
docker exec -it redis redis-cli set mykey "Hello Redis" ex 300
```

### 🔐 Authentication (if needed):
To enable password authentication, modify the docker-compose.yml command to:
```yaml
command: redis-server --appendonly yes --requirepass mypassword
```

Then connect with:
```bash
docker exec -it redis redis-cli -a mypassword
```

### 🔄 Backup Redis data:
```bash
docker exec -it redis redis-cli SAVE
docker cp redis:/data/dump.rdb ./redis-backup.rdb
```

### 📈 Monitor Redis commands in real-time:
```bash
docker exec -it redis redis-cli monitor
```