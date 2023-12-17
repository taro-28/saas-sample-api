include .env

.PHONY: gqlgenc

run:
	go run main.go
test:
	go test -v ./e2etest -count=1 -overlay=`testtime`
gqlgen:
	go get github.com/99designs/gqlgen@latest
	go run github.com/99designs/gqlgen generate
gqlgenc:
	go run e2etest/gqlgenc/main.go
xogen:
	xo schema -o models --src models/templates mysql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}/${DB_NAME}?tls=true
migrate:
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} ${DB_NAME} < ./db/schema.sql
migrate-dry:
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} ${DB_NAME} --dry-run < ./db/schema.sql