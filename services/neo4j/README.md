# 🔄 Neo4j

A complete Neo4j graph database setup for development and testing purposes.

## 🚀 Usage

### Using InfraCLI (Recommended)

```bash
# Start the service
infracli run neo4j

# Get connection information
infracli info neo4j

# Stop the service
infracli down neo4j
```

### Manual Usage

```bash
cd neo4j
docker-compose up -d
```

## ⚙️ Configuration

- **Browser Interface:** http://localhost:7474
- **Bolt Port:** 7687
- **Username:** neo4j
- **Password:** password
- **Connection URI:** bolt://localhost:7687
- **HTTPS Port:** 7473

## 💾 Data Persistence

Data is stored in named Docker volumes (`neo4j_data`, `neo4j_logs`, `neo4j_import`, and `neo4j_plugins`) which persist between container restarts. To reset the database, use `infracli down neo4j --volumes` or `docker-compose down -v` to remove the volumes.

## 🔧 Common Operations

### 💻 Connect using Cypher Shell:
```bash
docker exec -it neo4j cypher-shell -u neo4j -p password
```

### 📋 Create nodes and relationships:
```bash
docker exec -it neo4j cypher-shell -u neo4j -p password "CREATE (n:Person {name: 'John'}) RETURN n"
```

### 🔍 Query data:
```bash
docker exec -it neo4j cypher-shell -u neo4j -p password "MATCH (n:Person) RETURN n"
```

### 🗑️ Delete data:
```bash
docker exec -it neo4j cypher-shell -u neo4j -p password "MATCH (n) DETACH DELETE n"
```

### 🔄 Import data:
```bash
# Copy CSV file to import directory
docker cp mydata.csv neo4j:/var/lib/neo4j/import/

# Import using Cypher
docker exec -it neo4j cypher-shell -u neo4j -p password "LOAD CSV FROM 'file:///mydata.csv' AS row CREATE (:Person {name: row[0], age: toInteger(row[1])})"
```

### 📊 Export data:
```bash
docker exec -it neo4j cypher-shell -u neo4j -p password "MATCH (n:Person) RETURN n.name, n.age" > export.csv
```