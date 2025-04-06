.PHONY: launch_services launch_with_tests build_services stop up_db attach

launch_services: build_services
	docker compose up

launch_with_tests: build_services
	docker compose --profile test up

stop:
	docker compose down

build_services:
	docker compose build

up_db: 
	docker compose up postgres

attach:
	docker exec -it httpserver-postgres-1 psql -U user -d postgres_db