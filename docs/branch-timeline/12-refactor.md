# 12-refactor.md

## Added Input Validation with "github.com/go-playground/validator/v10"
```text
curl -s -X POST $BASE_URL/v1/customers   -H "Content-Type: application/json"   --data @test/customer.json
Validation error: Key: 'CreateCustomerRequest.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag
```

## Update phone number to e164 format

## Updated testing to use python for easier api testing

```text
# Running Customer Tests
python3 test/customer.py

✅ Created customer: 638

✅ Deleted customer 638: response 204
.
✅ First customer returned from http://localhost:8080/v1/customers/1: MARY
.
✅ Customer list is not empty from http://localhost:8080/v1/customers

✅ First customer Exists in list of customers
.
✅ API Health http://localhost:8080/health: ok
.
----------------------------------------------------------------------
Ran 4 tests in 0.022s

OK
# Finished Customer Tests...
```

```text
# Running Rental Tests...
python3 test/rental.py

✅ API /health: ok
.
✅ Rental created with ID: 16113
✅ Movie returned successfully (status 204)
.
✅ First rental created with ID: 16114
❌ Second rental failed as expected (status 500)
✅ Movie returned successfully (status 204)
.
----------------------------------------------------------------------
Ran 3 tests in 0.028s

OK
# Finished Rental Tests...
```

## Converted rest of tests
```
# Running Inventory Tests...
python3 test/inventory.py

📦 Testing: GET /v1/inventory
✅ Inventory list retrieved successfully
.
🎞️ Testing: GET /v1/inventory/available?film_id=1&store_id=2
✅ Available inventory retrieved successfully
.
🏬 Testing: GET /v1/inventory?store_id=1
✅ Store inventory retrieved successfully
.
----------------------------------------------------------------------
Ran 3 tests in 0.068s

OK
# Finished Inventory Tests...


# Running Store Tests...
python3 test/store.py

🏪 Testing: GET /v1/stores/1/inventory/summary
✅ Store inventory summary retrieved successfully
.
----------------------------------------------------------------------
Ran 1 test in 0.007s

OK
# Finished Store Tests...


# Running Film Tests...
python3 test/film.py

===== STARTING Film Tests =====

🎬 Testing: GET /v1/films
✅ Films list retrieved successfully
.
🎞️ Testing: GET /v1/films/1
✅ Film details retrieved successfully
.
🎭 Testing: GET /v1/films/1/with-actors-categories
✅ Enriched film data retrieved successfully
.
🔍 Testing: GET /v1/films/search?title=ACADEMY DINOSAUR
✅ Film search results retrieved successfully
.
----------------------------------------------------------------------
Ran 4 tests in 0.030s

OK
# Finished Film Tests...
```