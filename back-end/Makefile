include .env
.PHONY: test wire

dev:
	nodemon --watch .env --watch '**/*.go' --ext go,env --exec "go run main.go" --signal SIGTERM --ignore data/

wire:
	cd ./wire && go wire

lint:
	golangci-lint run --timeout 5m -v

uncommit:
	git reset HEAD^

##################
### Migrations ###
##################

install-migration:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

create-migration:
	migrate create -ext sql -dir ./migrations -seq ${NAME}

migrate:
	migrate -path ./migrations -database postgres://${BACKEND_DATABASE_ROOT_USERNAME}:${BACKEND_DATABASE_ROOT_PASSWORD}@${BACKEND_DATABASE_CONTAINER}:${BACKEND_DATABASE_PORT}/${BACKEND_DATABASE_NAME}?sslmode=disable up

migrate-down:
	migrate -path ./migrations -database "postgres://${BACKEND_DATABASE_ROOT_USERNAME}:${BACKEND_DATABASE_ROOT_PASSWORD}@${BACKEND_DATABASE_CONTAINER}:${BACKEND_DATABASE_PORT}/${BACKEND_DATABASE_NAME}?sslmode=disable" down ${N}

	