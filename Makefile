compose-up:
	docker compose up --build

compose-down:
	docker compose down -v

compose-service:
	docker compose up --build expenses-management -d