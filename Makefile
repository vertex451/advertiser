# section for everything
.PHONY: build_all
build_all:
	docker-compose build

.PHONY: run_all
run_all:
	docker-compose up -d

.PHONY: stop_all
stop_all:
	docker-compose down

.PHONY: restart_all
restart_all: stop_all run_all

# services only section
.PHONY: run_services
run_services:
	docker-compose up -d agency owner

.PHONY: stop_services
stop_services:
	docker-compose down agency owner
