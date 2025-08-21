compose-up:
	docker compose up --build

compose-down:
	docker compose down -v

compose-service:
	docker compose up --build expenses-management -d

gen_mock:
	cd repository/interface && mockgen -source=interface.go -package=_interface -destination=interface_mock.go