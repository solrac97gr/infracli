# 🐬 MySQL

A complete MySQL setup with persistence for development and testing purposes.

## 🚀 Usage

### Using InfraCLI (Recommended)

```bash
# Start the service
infracli run mysql

# Get connection information
infracli info mysql

# Stop the service
infracli down mysql
```

### Manual Usage

```bash
cd mysql
docker-compose up -d
```

## ⚙️ Configuration

- **Port:** 3306
- **Root Password:** rootpassword
- **Default Database:** testdb (created automatically)
- **Default User:** mysqluser
- **Default Password:** mysqlpassword
- **Connection String:** mysql://mysqluser:mysqlpassword@localhost:3306/testdb

## 💾 Data Persistence

Data is stored in a named Docker volume `mysql-data` which persists between container restarts. To reset the database, use `infracli down mysql --volumes` or `docker-compose down -v` to remove the volume as well.

## 📄 Adding Initial Database Scripts

You can place SQL scripts in the `init` directory to have them executed when the container is first created. Scripts are executed in alphabetical order.

Example:
```bash
# Create an init directory if it doesn't exist
mkdir -p init

# Add a sample script
echo "CREATE TABLE IF NOT EXISTS test_table (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255));" > init/01-create-tables.sql
echo "INSERT INTO test_table (name) VALUES ('test1'), ('test2');" > init/02-seed-data.sql
```

## 🔧 Common Operations

### 💻 Connect using the MySQL CLI client:
```bash
docker exec -it mysql mysql -u mysqluser -pmysqlpassword testdb
```

### 🔍 Execute SQL directly:
```bash
docker exec -it mysql mysql -u mysqluser -pmysqlpassword -e "SELECT * FROM testdb.test_table;"
```

### 📦 Backup a database:
```bash
docker exec mysql mysqldump -u mysqluser -pmysqlpassword testdb > backup.sql
```

### 🔄 Restore from backup:
```bash
cat backup.sql | docker exec -i mysql mysql -u mysqluser -pmysqlpassword testdb
```