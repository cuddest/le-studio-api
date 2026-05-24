#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-https://le-studio-api.onrender.com}"
ADMIN_EMAIL="${ADMIN_EMAIL:-}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-}"

echo "== Production smoke test =="
echo "BASE_URL=${BASE_URL}"

assert_200() {
  local path="$1"
  local status
  status=$(curl -sS -o /tmp/smoke_body.json -w "%{http_code}" "${BASE_URL}${path}")
  if [[ "$status" != "200" ]]; then
    echo "FAIL ${path} -> HTTP ${status}"
    cat /tmp/smoke_body.json || true
    exit 1
  fi
  echo "PASS ${path}"
}

assert_200 "/api/v1/coaches"
assert_200 "/api/v1/training-types"
assert_200 "/api/v1/pack-templates"
assert_200 "/api/v1/schedules"

if [[ -n "${ADMIN_EMAIL}" && -n "${ADMIN_PASSWORD}" ]]; then
  echo "Running protected admin smoke checks..."
  login_json=$(curl -sS -X POST "${BASE_URL}/api/v1/admin/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}")

  admin_token=$(echo "$login_json" | jq -r '.data.access_token // .access_token // empty')
  if [[ -z "$admin_token" ]]; then
    echo "FAIL admin login"
    echo "$login_json"
    exit 1
  fi

  me_status=$(curl -sS -o /tmp/admin_me.json -w "%{http_code}" \
    -H "Authorization: Bearer ${admin_token}" \
    "${BASE_URL}/api/v1/admin/auth/me")

  if [[ "$me_status" != "200" ]]; then
    echo "FAIL /api/v1/admin/auth/me -> HTTP ${me_status}"
    cat /tmp/admin_me.json || true
    exit 1
  fi

  echo "PASS /api/v1/admin/auth/me"
else
  echo "Skipping protected checks (ADMIN_EMAIL/ADMIN_PASSWORD not set)."
fi

echo "Smoke test passed."