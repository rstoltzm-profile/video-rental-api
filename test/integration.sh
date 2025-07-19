#!/usr/bin/env bash

set -euo pipefail

BASE_URL="http://localhost:8080"

echo "🔍 Running integration tests... requires jq"
echo "## Server Tests"
# Test 1: Health check
echo -n "✅ /health ... "
curl -s $BASE_URL/health | grep -q '"status":"ok"' && echo "OK" || (echo "FAIL" && exit 1)

##########################################################
echo "## Customer Tests"
# Test 1: GET /v1/customers (expects non-empty list)
echo -n "✅ /v1/customers ... "
RESPONSE=$(curl -s $BASE_URL/v1/customers)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty customer list"
    exit 1
fi

# Test 2: GET /v1/customers/1 (expects ID = 1)
echo -n "✅ /v1/customers/1 ... "
RESPONSE=$(curl -s $BASE_URL/v1/customers/1)
if echo "$RESPONSE" | jq -e '.id == 1' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Customer ID mismatch"
    exit 1
fi

# Test 3: Add customer /v1/customers
echo -n "✅ POST /v1/customers ... "

CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/v1/customers \
  -H "Content-Type: application/json" \
  --data @test/customer.json)

CUSTOMER_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

if [[ "$CUSTOMER_ID" =~ ^[0-9]+$ ]]; then
  echo "OK (created customer ID = $CUSTOMER_ID)"
else
  echo "FAIL - could not extract customer ID"
  echo "$CREATE_RESPONSE"
  exit 1
fi

# Test 4: Delete /v1/customers/1
echo -n "✅ DELETE /v1/customers/$CUSTOMER_ID ... "

DELETE_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE $BASE_URL/v1/customers/$CUSTOMER_ID)

if [ "$DELETE_STATUS" -eq 204 ]; then
  echo "OK"
else
  echo "FAIL - unexpected status code: $DELETE_STATUS"
  exit 1
fi

##########################################################
echo "## Rental Tests"
# Test 1: GET /v1/rentals (expects non-empty list)
echo -n "✅ /v1/rentals ... "
RESPONSE=$(curl -s $BASE_URL/v1/rentals)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty rentals list"
    exit 1
fi

# Test 2: GET /v1/rentals?late=true (expects non-empty list)
echo -n "✅ /v1/rentals?late=true ... "
RESPONSE=$(curl -s $BASE_URL/v1/rentals?late=true)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty rentals list"
    exit 1
fi

# Test 3: GET /v1/rentals?customer_id=1 (expects non-empty list)
echo -n "✅ /v1/rentals?customer_id=373 ... "
RESPONSE=$(curl -s $BASE_URL/v1/rentals?customer_id=373)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty rentals list"
    exit 1
fi

# Test 4: GET /v1/rentals?customer_id=1&late=true (expects non-empty list)
echo -n "✅ /v1/rentals?customer_id=1&late=true ... "
RESPONSE=$(curl -s $BASE_URL/v1/rentals?customer_id=373&late=true )
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty rentals list"
    exit 1
fi

##########################################################
echo "## Inventory Tests"
# Test 1: GET /v1/inventory (expects non-empty list)
echo -n "✅ /v1/inventory ... "
RESPONSE=$(curl -s $BASE_URL/v1/inventory)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty inventory list"
    exit 1
fi

# Test 2: GET /v1/inventory with store_id (expects non-empty list)
echo -n "✅ /v1/inventory?store_id=1 ... "
RESPONSE=$(curl -s $BASE_URL/v1/inventory?store_id=1)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty inventory list"
    exit 1
fi

echo "## Store Tests"
# Test 1: GET /v1/stores/1/inventory/summary (expects non-empty list)
echo -n "✅ /v1/stores/1/inventory/summary ... "
RESPONSE=$(curl -s $BASE_URL/v1/stores/1/inventory/summary)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty inventory list"
    exit 1
fi

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

# Test 4: GET /v1/films/1/with-actors(expects non-empty list)
echo -n "✅ /v1/films/1/with-actors ... "
RESPONSE=$(curl -s $BASE_URL/v1/films/1/with-actors-categories)
if echo "$RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo "OK"
else
    echo "FAIL - Expected non-empty films list"
    exit 1
fi
