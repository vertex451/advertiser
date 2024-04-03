.PHONY: go_run_owner
go_run_owner:
	cd ./golang/owner && go run cmd/main.go

.PHONY: tests
tests: integration-test-up
	go test -v -cover ./...
	make integration-test-down

.PHONY: integration-test-up
integration-test-up:
	docker-compose -f integration_tests/docker-compose-owner-test.yml up -d

.PHONY: integration-test-down
integration-test-down:
	docker-compose -f integration_tests/docker-compose-owner-test.yml down