.PHONY: launch_services launch_with_tests build_services stop_services

launch_services: build_services
	docker compose up

launch_with_tests: build_services
	docker compose --profile test up --abort-on-container-exit --exit-code-from server_test

stop_services:
	docker compose down

build_services:
	docker compose build