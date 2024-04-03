include ./golang/agency/Makefile
include ./golang/owner/Makefile

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


# postgres only section
.PHONY: run_postgres
run_postgres:
	docker-compose up -d postgres

# services only section
.PHONY: build_services
build_services:
	docker-compose build agency owner

.PHONY: push_services
push_services:
	docker push vertex451/agency:latest && docker push vertex451/owner:latest

.PHONY: run_services
run_services:
	docker-compose up -d agency owner

.PHONY: stop_services
stop_services:
	docker-compose down agency owner


.PHONY: build_agency
build_agency:
	docker-compose build agency


.PHONY: run_agency
run_agency:
	docker-compose up -d agency