.DEFAULT_GOAL := help

DOCKER_DB_CONTAINER = postgres_db

DOCKER_DB_CONTAINERS = -f docker/db.yml
DOCKER_DB_CONTAINERS += -f docker/pgAdmin.yml

# use bash as shell for command execution
ifeq ($(UNAME),Darwin)
	SHELL := /opt/local/bin/bash
	OS_X  := true
	export SHELL:=/opt/local/bin/bash
	export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit
else
	OS_DEB  := true
	SHELL := /bin/bash
	export SHELL:=/bin/bash
	export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit
endif

.ONESHELL:


help:
	@echo "Hello from help"
	@echo " "	
	@echo "Make cmd's "
	@echo "- stop                               stop docker db containers"
	@echo "- docker                             build and run postgres db and pgAdmin cotainer"
	@echo " "


stop:
	- docker compose $(DOCKER_DB_CONTAINERS) -p test_backend down -t 2


.PHONY: docker

docker:
	function tearDown {
			make stop
	}
	trap tearDown EXIT
	- docker compose $(DOCKER_DB_CONTAINERS) -p test_backend up

run_app:
	- mvn spring-boot:run
