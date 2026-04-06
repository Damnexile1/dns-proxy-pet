#!/bin/bash

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker compose exec -T postgres pg_isready -U dnsproxy; do
  sleep 1
done

echo "PostgreSQL is ready!"

# Apply migrations
echo "Applying migrations..."

# Run migrations directly using psql
for file in migrations/*.up.sql; do
  echo "Applying migration: $file"
  docker compose exec -T postgres psql -U dnsproxy -d dnsproxy < "$file"
done

echo "Migrations applied successfully!"
