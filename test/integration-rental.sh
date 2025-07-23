#!/usr/bin/env bash

set -euo pipefail

BASE_URL_2="http://localhost:8080"

##########################################################
echo "## Rent a movie"
# Test 1: Rent movie
echo -n "✅ POST /v1/rentals ... "

CREATE_RESPONSE=$(curl -s -X POST $BASE_URL_2/v1/rentals \
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

RETURN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL_2/v1/rentals/$RENTAL_ID/return)

if [ "$RETURN_STATUS" -eq 204 ]; then
  echo "OK"
else
  echo "FAIL - unexpected status code: $RETURN_STATUS"
  exit 1
fi

##########################################################
# Test 3: Rent a movie that's already rented
echo -n "✅ POST /v1/rentals (first attempt) ... "

CREATE_RESPONSE_1=$(curl -s -X POST $BASE_URL_2/v1/rentals \
  -H "Content-Type: application/json" \
  --data @test/rental.json)

RENTAL_ID_1=$(echo "$CREATE_RESPONSE_1" | jq -r '.id')

if [[ "$RENTAL_ID_1" =~ ^[0-9]+$ ]]; then
  echo "OK (created rental ID = $RENTAL_ID_1)"
else
  echo "FAIL - could not extract rental ID"
  echo "$CREATE_RESPONSE_1"
  exit 1
fi

echo -n "✅ POST /v1/rentals (second attempt, should fail) ... "

CREATE_RESPONSE_2=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL_2/v1/rentals \
  -H "Content-Type: application/json" \
  --data @test/rental.json)

RESPONSE_BODY=$(echo "$CREATE_RESPONSE_2" | head -n 1)
HTTP_STATUS=$(echo "$CREATE_RESPONSE_2" | tail -n 1)

if [[ "$HTTP_STATUS" -ge 400 ]]; then
  echo "OK (duplicate rental rejected with status $HTTP_STATUS)"
else
  echo "FAIL - duplicate rental was allowed"
  echo "Response: $RESPONSE_BODY"
fi

RETURN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL_2/v1/rentals/$RENTAL_ID_1/return)