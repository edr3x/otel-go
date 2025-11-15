BASE_URL ?= http://localhost:8080

request-success:
	@echo "=> Successful request"
	@curl -s -w "\n" "$(BASE_URL)/users/123323"

request-failed:
	@echo "=> Request should fail (no post, no asset)"
	@curl -s -w "\n" "$(BASE_URL)/users/9999"

request-partial:
	@echo "=> Partial request (omit asset service)"
	@curl -s -w "\n" "$(BASE_URL)/users/1234"

.PHONY: request-success request-failed request-partial
