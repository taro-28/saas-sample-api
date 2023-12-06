include .env

run:
	go run main.go
gqlgen:
	go get github.com/99designs/gqlgen@latest
	go run github.com/99designs/gqlgen generate
xogen:
	xo schema "${DSN_FOR_XO}" -o models