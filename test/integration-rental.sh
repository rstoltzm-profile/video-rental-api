#!/usr/bin/env bash

set -euo pipefail

BASE_URL="http://localhost:8080"

##########################################################
echo "## Rent a movie"
# Test 1: Rent movie
echo -n "✅ POST /v1/rentals ... "

CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/v1/rentals \
  -H "Content-Type: application/json" \
  --data @test/rental.json)

RENTAL_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

if [[ "$RENTAL_ID" =~ ^[0-9]+$ ]]; then
  echo "OK (created rental ID = $RENTAL_ID)"
else
  echo "FAIL - could not extract rental ID"
  echo "$CREATE_RESPONSE"
  exit 1
fi

# Test 2: Return movie
echo -n "✅ POST /v1/rentals/$RENTAL_ID/return ... "

RETURN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/v1/rentals/$RENTAL_ID/return)

if [ "$RETURN_STATUS" -eq 204 ]; then
  echo "OK"
else
  echo "FAIL - unexpected status code: $RETURN_STATUS"
  exit 1
fi