# 12-refactor.md

## Added Input Validation with "github.com/go-playground/validator/v10"
```text
curl -s -X POST $BASE_URL/v1/customers   -H "Content-Type: application/json"   --data @test/customer.json
Validation error: Key: 'CreateCustomerRequest.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag
```

## Update phone number to e164 format