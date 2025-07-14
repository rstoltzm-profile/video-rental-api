# 4 - Extend API for Customer

✅ Summary:
- Add logic for getting address_id, city_id, insert address, insert customer
- Add Create Customer
- Add Delete Customer


## integration-test
```bash
make integration-test
```

```
bash test/integration.sh
🔍 Running integration tests... requires jq
✅ /health ... OK
✅ /v1/customers ... OK
✅ /v1/customers/1 ... OK
✅ POST /v1/customers ... OK (created customer ID = 606)
✅ DELETE /v1/customers/606 ... OK
```