CURRENT_DIR = $(shell pwd)

DB_URL := postgres://postgres:123321@localhost:5432/tender_management?sslmode=disable

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

mig-up:
	migrate -path migrations -database '${DB_URL}' -verbose up

mig-down:
	migrate -path migrations -database '${DB_URL}' -verbose down

mig-force:
	migrate -path migrations -database '${DB_URL}' -verbose force 1

mig-create:
	migrate create -ext sql -dir migrations -seq tender_management

run_db:
	docker compose build db && docker compose up -d db && docker compose up -d redis &&make mig-up

run:
	docker compose build app && docker compose up -d app

swag-gen:
	~/go/bin/swag init -g internal/controller/http/router.go -o docs
#   rm -r db/migrations

