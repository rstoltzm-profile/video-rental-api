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