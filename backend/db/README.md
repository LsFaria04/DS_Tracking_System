# Database Setup

This directory contains database schemas and migrations for the order tracking system.

## Structure

- `migrations/` - SQL migration files for schema changes
- `seeds/` - Sample data for development

## Running Migrations

To apply the initial schema:

```bash
# Start the database
docker compose -f compose.dev.yml up postgres -d

# Apply migration
docker compose -f compose.dev.yml exec postgres psql -U admin -d postgres -f /docker-entrypoint-initdb.d/001_init_schema.sql

# (Optional) Load seed data
docker compose -f compose.dev.yml exec postgres psql -U admin -d postgres -f /docker-entrypoint-initdb.d/dev_data.sql
```

## Schema Overview

### Tables

**orders**
- Stores basic order information
- Tracks current status
- Links to customer

**order_status_history**
- Immutable audit log of all status changes
- Provides transparency and traceability
- Timestamps all transitions

### Status Lifecycle

1. `IN_PRODUCTION` - Order is being manufactured
2. `READY_FOR_SHIPMENT` - Production complete, awaiting pickup
3. `SHIPPED` - Package in transit
4. `RECEIVED` - Delivered to customer
