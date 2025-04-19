# 🐘 PostgreSQL

A PostgreSQL setup with initialization script for development and testing.

## 🚀 Usage

### Using InfraCLI (Recommended)

```bash
# Start the service
infracli run postgres

# Get connection information
infracli info postgres

# Stop the service
infracli down postgres
```

### Manual Usage

```bash
cd postgres
docker-compose up -d
```

## ⚙️ Configuration

- **Port:** 5432
- **Username:** postgres
- **Password:** postgres
- **Default Database:** testdb (created via initialization script)
- **Connection String:** postgresql://postgres:postgres@localhost:5432/testdb

## 💾 Data Persistence

Data is stored in a named Docker volume `postgres_data` which persists between container restarts. To reset the database, use `infracli down postgres --volumes` or `docker-compose down -v` to remove the volume as well.

## 📄 Initialization Scripts

You can customize the `initdb.sql` file in the postgres directory to run SQL commands when the container first starts. By default, this creates a test database.

## 🔧 Common Operations

### 💻 Connect using psql:
```bash
docker exec -it db psql -U postgres
```

### 🔍 Connect to a specific database:
```bash
docker exec -it db psql -U postgres -d testdb
```

### 📋 List databases:
```bash
docker exec -it db psql -U postgres -c "\l"
```

### 📦 Backup a database:
```bash
docker exec -it db pg_dump -U postgres testdb > backup.sql
```

### 🔄 Restore from backup:
```bash
cat backup.sql | docker exec -i db psql -U postgres -d testdb
```

### 📄 Execute SQL script:
```bash
docker exec -i db psql -U postgres -d testdb < script.sql
```

### 📊 View tables in a database:
```bash
docker exec -it db psql -U postgres -d testdb -c "\dt"
```