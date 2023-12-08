include .env

run:
	go run main.go
gqlgen:
	go get github.com/99designs/gqlgen@latest
	go run github.com/99designs/gqlgen generate
xogen:
	xo schema -o models --src models/templates mysql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}/${DB_NAME}?tls=true
migrate:
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} ${DB_NAME} < ./db/schema.sql
migrate-dry:
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} ${DB_NAME} --dry-run < ./db/schema.sql