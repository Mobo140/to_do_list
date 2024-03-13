build:	
	docker-compose --build to_do_list

run:
	docker-compose up to_do_list

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up

swag:
	swag init -g cmd/main.go