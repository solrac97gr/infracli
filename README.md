# 🚀 Development Infrastructure

This repository contains Docker Compose configurations for various services used for testing and development environments.

## 🏁 Quick Start

### 🔧 Using the InfraCLI Tool

You can use our CLI tool to manage infrastructure services easily:

```bash
# Install the CLI tool
cd infracli
./scripts/install.sh

# List available services
infracli list

# Get connection information for a service
infracli info mysql

# Start services
infracli run mysql
infracli run mongo elasticsearch-kibana
infracli run all

# Stop services
infracli down mysql
infracli down all
```

## 🧰 Available Infrastructure

Each service has its own detailed documentation available in its directory.

### 🍃 MongoDB

A simple MongoDB setup for development and testing purposes.

[📄 View MongoDB Documentation](services/mongo/README.md)

### 🐘 PostgreSQL

A PostgreSQL setup with initialization script for development and testing.

[📄 View PostgreSQL Documentation](services/postgres/README.md)

### 🔴 Redis

A Redis setup with persistence for development and caching needs.

[📄 View Redis Documentation](services/redis/README.md)

### 🐬 MySQL

A complete MySQL setup with data persistence and initialization scripts.

[📄 View MySQL Documentation](services/mysql/README.md)

### 🔍 Elasticsearch and Kibana

A complete setup for Elasticsearch and Kibana with data persistence.

[📄 View Elasticsearch & Kibana Documentation](services/elasticsearch-kibana/README.md)

## 🔄 General Usage

Our recommended workflow is to use the InfraCLI tool for all operations:

1. Install the InfraCLI tool (see below)
2. Use `infracli run` to start services
3. Use `infracli info` to get connection details
4. Use `infracli down` to stop services when finished

## 💾 Data Persistence

All databases are configured with named volumes to persist data between container restarts. If you need to reset the database, use `infracli down [service] --volumes` to remove the volumes as well.

## ➕ Adding New Infrastructure

Follow the pattern established in the existing directories to add configurations for additional databases as needed:

1. Create a directory for your service in `services/`
2. Add a `docker-compose.yml` file
3. Include a README.md with detailed usage information
4. The CLI tool will automatically detect your new service

## 🛠️ InfraCLI Tool

The repository includes a Go-based CLI tool for efficient service management. Learn more in the [InfraCLI README](cli/README.md).

### ✨ Key Features

- Start and stop services with simple commands
- Get connection details for any service
- Manage multiple services at once
- Automatic service discovery
- Proper cleanup with volume removal option

### 📦 Installation

```bash
cd infracli
./scripts/install.sh
```

Once installed, you can use `infracli` from anywhere in your system to manage the infrastructure services.
