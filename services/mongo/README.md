# 🍃 MongoDB

A simple MongoDB setup for development and testing purposes.

## 🚀 Usage

### Using InfraCLI (Recommended)

```bash
# Start the service
infracli run mongo

# Get connection information
infracli info mongo

# Stop the service
infracli down mongo
```

### Manual Usage

```bash
cd mongo
docker-compose up -d
```

## ⚙️ Configuration

- **Port:** 27017
- **Username:** mongo
- **Password:** password
- **Authentication Database:** admin
- **Connection String:** mongodb://mongo:password@localhost:27017/admin

## 💾 Data Persistence

Data is stored in a named Docker volume `mongodb_data` which persists between container restarts. To reset the database, use `infracli down mongo --volumes` or `docker-compose down -v` to remove the volume as well.

## 🔧 Common Operations

### 💻 Connect using the MongoDB Shell:
```bash
docker exec -it mongo mongosh --username mongo --password password
```

### 📦 Create a new database:
```bash
docker exec -it mongo mongosh --username mongo --password password --eval "use myNewDB; db.createCollection('myCollection')"
```

### 📋 List databases:
```bash
docker exec -it mongo mongosh --username mongo --password password --eval "show dbs"
```

### 🔍 Query data:
```bash
docker exec -it mongo mongosh --username mongo --password password --eval "use myDB; db.myCollection.find()"
```

### 📤 Export data:
```bash
docker exec -it mongo mongoexport --username mongo --password password --db myDB --collection myCollection --out /data/myCollection.json
```

### 📥 Import data:
```bash
# Copy file into container first
docker cp myCollection.json mongo:/tmp/
docker exec -it mongo mongoimport --username mongo --password password --db myDB --collection myCollection --file /tmp/myCollection.json
```