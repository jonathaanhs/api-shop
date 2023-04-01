
test:
	go clean -testcache
	go test ./... --cover

run-pg:
	docker-compose -f ./deploy/pg.yaml up --build -d

stop-pg:
	docker-compose -f ./deploy/pg.yaml down

migrate-up:
	migrate -database "postgresql://dbuser:dbpass@:5432/dbname?sslmode=disable" -path ./database/pg/migration up

migrate-down:
	migrate -database "postgresql://dbuser:dbpass@:5432/dbname?sslmode=disable" -path ./database/pg/migration down