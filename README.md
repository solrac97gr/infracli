# Development Infrastructure

This repository contains Docker Compose configurations for various services used for testing and development environments.

## Quick Start

### Using the CLI Tool (Recommended)

You can use our CLI tool to manage infrastructure services:

```bash
# Install the CLI tool
cd infracli
./scripts/install.sh

# Start services
infracli run mysql
infracli run mongo elasticsearch-kibana
infracli run all

# Stop services
infracli down mysql
infracli down all

# List available services
infracli list
```

### Using the Shell Script

Alternatively, you can use the provided shell script:

```bash
# Start a single service
./run.sh mysql

# Start multiple services
./run.sh mongo elasticsearch-kibana

# Start all available services
./run.sh all
```

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

### Using the run.sh Script

The `run.sh` script provides a convenient way to start services:

```bash
# Display help information
./run.sh

# Start one or more services
./run.sh [service1] [service2] ... [serviceN]
```

To stop services, you'll still need to use the docker-compose command in each service directory:

```bash
cd [service-directory]
docker-compose down
```

A future enhancement could include a stop capability in the script.

## Data Persistence

All databases are configured with named volumes to persist data between container restarts. If you need to reset the database, use `docker-compose down -v` to remove the volumes as well.

## Adding New Infrastructure

Follow the pattern established in the existing directories to add configurations for additional databases as needed.

## InfraCLI Tool

The repository includes a Go-based CLI tool for more efficient service management. Learn more in the [InfraCLI README](infracli/README.md).

### Key Features

- Start and stop services with simple commands
- Manage multiple services at once
- Automatic service discovery
- Proper cleanup with volume removal option

### Installation

```bash
cd infracli
./scripts/install.sh
```

Once installed, you can use `infracli` from anywhere in your system to manage the infrastructure services.
