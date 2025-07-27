# 0-post-mvp-validation.md

## Notes
- I ended up adding more tests to customers service.go
- I had an issue with mocking transactions being annoying
- I decided to refactor customer service to match rental service type method of ctx and txn
- This made it easier to write a test for service layer and more consistent with rest of code
- 