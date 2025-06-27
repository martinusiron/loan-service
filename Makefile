docker-rebuild:
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up

swagger:
	rm -rf docs/
	swag init -g cmd/main.go --parseDependency --parseInternal

test:
	go test -v ./tests/...

test-integration:
	DB_HOST=localhost DB_NAME=loan_db go test -v ./tests/...

	