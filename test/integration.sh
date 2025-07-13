#!/usr/bin/env bash

set -euo pipefail

BASE_URL="http://localhost:8080"

echo "🔍 Running integration tests... requires jq"

# Test 1: Health check
echo -n "✅ /health ... "
curl -s $BASE_URL/health | grep -q '"status":"ok"' && echo "OK" || (echo "FAIL" && exit 1)

# Test 2: GET /v1/customers (expects non-empty list)
echo -n "✅ /v1/customers ... "
RESPONSE=$(curl -s $BASE_URL/v1/customers)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty customer list"
    exit 1
fi

# Test 3: GET /v1/customers/1 (expects ID = 1)
echo -n "✅ /v1/customers/1 ... "
RESPONSE=$(curl -s $BASE_URL/v1/customers/1)
if echo "$RESPONSE" | jq -e '.id == 1' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Customer ID mismatch"
    exit 1
fi

echo "🎉 All integration tests passed!"
