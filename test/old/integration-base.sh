#!/usr/bin/env bash

set -euo pipefail

BASE_URL="http://localhost:8080"

echo "🔍 Running integration tests... requires jq"
echo "## Server Tests"
# Test 1: Health check
echo -n "✅ /health ... "
curl -s $BASE_URL/health | grep -q '"status":"ok"' && echo "OK" || (echo "FAIL" && exit 1)


##########################################################
echo "## Film Tests"
# Test 1: GET /v1/films (expects non-empty list)
echo -n "✅ /v1/films ... "
RESPONSE=$(curl -s $BASE_URL/v1/films)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty films list"
    exit 1
fi

# Test 2: GET /v1/films/1 (expects non-empty list)
echo -n "✅ /v1/films/1 ... "
RESPONSE=$(curl -s $BASE_URL/v1/films/1)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty films list"
    exit 1
fi

# Test 3: GET /v1/films/search?title=ACADEMY DINOSAUR (expects non-empty list)
echo -n "✅ /v1/films/search?title=ACADEMY DINOSAUR ... "
RESPONSE=$(curl -s "$BASE_URL/v1/films/search?title=ACADEMY%20DINOSAUR")
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty films list"
    exit 1
fi

# Test 4: GET /v1/films/1/with-actors-categories (expects non-empty list)
echo -n "✅ /v1/films/1/with-actors-categories ... "
RESPONSE=$(curl -s $BASE_URL/v1/films/1/with-actors-categories)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty films list"
    exit 1
fi

echo "## Done with base tests"
echo ""