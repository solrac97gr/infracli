# Development Infrastructure

This repository contains Docker Compose configurations for various services used for testing and development environments.

## Available Infrastructure

### MongoDB

A simple MongoDB setup for development and testing purposes.

#### Usage

```bash
cd mongo
docker-compose up -d
```

#### Configuration

- **Port:** 27017
- **Username:** mongo
- **Password:** password
- **Connection String:** mongodb://mongo:password@localhost:27017

### PostgreSQL

A PostgreSQL setup with initialization script for development and testing.

#### Usage

```bash
cd postgres
docker-compose up -d
```

#### Configuration

- **Port:** 5432
- **Username:** postgres
- **Password:** postgres
- **Default Database:** testdb (created automatically)
- **Connection String:** postgresql://postgres:postgres@localhost:5432/testdb

### MySQL

A complete MySQL setup with data persistence and initialization scripts.

#### Usage

```bash
cd mysql
docker-compose up -d
```

#### Configuration

- **Port:** 3306
- **Root Password:** rootpassword
- **Default Database:** testdb
- **Default User:** mysqluser
- **Default Password:** mysqlpassword
- **Connection String:** mysql://mysqluser:mysqlpassword@localhost:3306/testdb

### Elasticsearch and Kibana

A complete setup for Elasticsearch and Kibana with data persistence.

#### Usage

```bash
cd elasticsearch-kibana
docker-compose up -d
```

#### Configuration

##### Elasticsearch
- **Port:** 9200
- **Internal Port:** 9300
- **Security:** Disabled for development
- **Connection URL:** http://localhost:9200

##### Kibana
- **Port:** 5601
- **Dashboard URL:** http://localhost:5601

## General Usage

1. Navigate to the directory of the database you need
2. Start the container: `docker-compose up -d`
3. Connect to the database using the provided connection details
4. Stop the container when finished: `docker-compose down`

## Data Persistence

All databases are configured with named volumes to persist data between container restarts. If you need to reset the database, use `docker-compose down -v` to remove the volumes as well.

## Adding New Infrastructure

Follow the pattern established in the existing directories to add configurations for additional databases as needed.
