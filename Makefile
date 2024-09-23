run_docker:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

start_docker:
	docker start postgres

stop_docker:
	docker stop postgres

create_db:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

drop_db:
	docker exec -it postgres dropdb simple_bank

up_migrate:
	migrate -database postgres://root:secret@localhost:5432/simple_bank?sslmode=disable -path db/migrations -verbose up

down_migrate:
	migrate -database postgres://root:secret@localhost:5432/simple_bank?sslmode=disable -path db/migrations -verbose down

sqlc_generate:
	sqlc generate
	
test:
	go test -timeout 30s -cover ./... -v

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/longln/simplebank/db/sqlc Store

1up_migrate:
	migrate -database postgres://root:secret@localhost:5432/simple_bank?sslmode=disable -path db/migrations -verbose up 1

1down_migrate:
	migrate -database postgres://root:secret@localhost:5432/simple_bank?sslmode=disable -path db/migrations -verbose down 1



.PHONY: 1up_migrate 1down_migrate mockgen server test sqlc_generate run_docker start_docker stop_docker create_db drop_db up_migrate down_migrate
