include .env

run:
	go run main.go
gqlgen:
	go get github.com/99designs/gqlgen@v0.17.40
	go run github.com/99designs/gqlgen generate
xogen:
	xo schema "${DSN_FOR_XO}" -o models