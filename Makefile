compose-up:
	docker compose up --build

compose-down:
	docker compose down -v

compose-service:
	docker compose up --build expenses-management -d

gen_mock:
	cd repository/interface && mockgen -source=interface.go -package=_interface -destination=interface_mock.go

test:
	go test ./service/... -coverprofile=coverage.out

test-coverage:
	go test ./service/... -coverprofile=coverage.out
	go tool cover -html=coverage.out
