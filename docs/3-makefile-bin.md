# 1 - Makefile binaries

✅ Summary:
- Add basic Makefile
- Add make integration-test

## integration-test
```bash
make integration-test
```

## integration-test output
```text
bash test/integration.sh
🔍 Running integration tests... requires jq
✅ /health ... OK
✅ /v1/customers ... OK
✅ /v1/customers/1 ... OK
🎉 All integration tests passed!
```