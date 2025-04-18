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

## General Usage

1. Navigate to the directory of the database you need
2. Start the container: `docker-compose up -d`
3. Connect to the database using the provided connection details
4. Stop the container when finished: `docker-compose down`

## Data Persistence

All databases are configured with named volumes to persist data between container restarts. If you need to reset the database, use `docker-compose down -v` to remove the volumes as well.

## Adding New Infrastructure

Follow the pattern established in the existing directories to add configurations for additional databases as needed.
