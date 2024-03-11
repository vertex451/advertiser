.PHONY: tests
tests: integration-test-up
	go test -v -cover ./...
	make integration-test-down

.PHONY: integration-test-up
integration-test-up:
	docker-compose -f integration_tests/docker-compose-test.yml up -d

.PHONY: integration-test-down
integration-test-down:
	docker-compose -f integration_tests/docker-compose-test.yml down