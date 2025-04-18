# Elasticsearch and Kibana

A complete setup for Elasticsearch and Kibana with persistence for development and testing purposes.

## Usage

```bash
cd elasticsearch-kibana
docker-compose up -d
```

## Configuration

### Elasticsearch
- **Port:** 9200
- **Internal Port:** 9300
- **Security:** Disabled for development (xpack.security.enabled=false)
- **Java Heap Size:** 512MB min/max
- **Connection URL:** http://localhost:9200

### Kibana
- **Port:** 5601
- **Dashboard URL:** http://localhost:5601
- **Elasticsearch Connection:** Preconfigured to http://elasticsearch:9200

## Data Persistence

Data is stored in a named Docker volume `elasticsearch-data` which persists between container restarts. To reset the database, use `docker-compose down -v` to remove the volume as well.

## Common Operations

### Check Elasticsearch status:
```bash
curl http://localhost:9200
```

### View indices:
```bash
curl http://localhost:9200/_cat/indices
```

### Create an index:
```bash
curl -X PUT "localhost:9200/my-index"
```

### Basic authentication (when enabled):
If you enable security by changing `xpack.security.enabled=true`, you'll need to set up passwords and use:
```bash
curl -u elastic:yourpassword http://localhost:9200
```