#!/usr/bin/env bash
set -euo pipefail

# Usage:
# ADMIN_SETUP_KEY=your_key API_BASE_URL=http://localhost:8080 \
#   ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=secret ADMIN_NAME="Admin User" \
#   ./scripts/seed_first_admin.sh

API_BASE_URL=${API_BASE_URL:-http://localhost:8080}
ADMIN_SETUP_KEY=${ADMIN_SETUP_KEY:-}
ADMIN_EMAIL=${ADMIN_EMAIL:-admin@example.com}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-password}
ADMIN_NAME=${ADMIN_NAME:-Admin}

if [ -z "$ADMIN_SETUP_KEY" ]; then
  echo "ERROR: ADMIN_SETUP_KEY must be set in environment to call setup endpoint."
  echo "Example: ADMIN_SETUP_KEY=your_key API_BASE_URL=http://localhost:8080 ./scripts/seed_first_admin.sh"
  exit 1
fi

echo "Seeding first admin to $API_BASE_URL using setup key"

payload=$(jq -n --arg email "$ADMIN_EMAIL" --arg password "$ADMIN_PASSWORD" --arg name "$ADMIN_NAME" '{email: $email, password: $password, name: $name}')

if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR: 'jq' is required to run this script. Install it and retry."
  exit 1
fi

resp=$(curl -sS -w "\nHTTP_STATUS:%{http_code}" -X POST "$API_BASE_URL/api/v1/admin/auth/register" \
  -H "Content-Type: application/json" \
  -H "x-admin-setup-key: $ADMIN_SETUP_KEY" \
  -d "$payload")

echo "$resp"

http_status=$(echo "$resp" | tr '\n' ' ' | sed -n 's/.*HTTP_STATUS:\([0-9][0-9][0-9]\)$/\1/p')

if [ "$http_status" = "201" ] || [ "$http_status" = "200" ]; then
  echo "Admin created successfully."
  exit 0
else
  echo "Failed to create admin (status: $http_status). See response above."
  exit 2
fi
