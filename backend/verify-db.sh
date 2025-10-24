#!/bin/bash

# Verify database tables and data

echo "Tables:"
docker compose -f compose.dev.yml exec -T postgres psql -U admin -d postgres -c "\dt"

echo ""
echo "Orders:"
docker compose -f compose.dev.yml exec -T postgres psql -U admin -d postgres -c "SELECT * FROM orders;"
