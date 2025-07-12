# Cookbook Deployment Guide

## Local Development

For local development with locally built images:

```bash
# Build and start with local image (default behavior)
docker-compose up --build

# Or just start if already built
docker-compose up
```

This uses `docker-compose.override.yml` automatically to build the image locally.

```bash
# Here's a one-liner
docker compose down && docker buildx build -t cookbook-app:latest -f Dockerfile.cookbook-app . && docker compose up -d
```

## Production Deployment

For production using GitHub Container Registry:

```bash
# Using production override file
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Or set environment variable
export COOKBOOK_IMAGE=ghcr.io/mryachanin/cookbook-app:latest
docker-compose up -d
```

## Environment Variables

- `COOKBOOK_IMAGE`: Docker image to use (defaults to GitHub registry)
- `COUCHDB_USER`: CouchDB username (defaults to "admin")
- `COUCHDB_PASSWORD`: CouchDB password (defaults to "admin")

## Database Setup

After starting the containers, set up the database:

```bash
# Build database tools
go build -o cookbook-db ./cmd/db

# Set up database
./cookbook-db -o setup

# Drop the database
./cookbook-db -o drop

# Import sample data
./cookbook-db -o import -i ./testdata

# Export all data into a folder
./cookbook-db -o export -e ./folder-to-export-into
```