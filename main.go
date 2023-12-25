package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/taro-28/saas-sample-api/db"
	"github.com/taro-28/saas-sample-api/gql"
	loaders "github.com/taro-28/saas-sample-api/loader"
)

const defaultPort = "8080"

func main() {
	// Load connection string from .env file
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, cleanup, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer cleanup()
	var srv http.Handler = handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{
		DB: db,
	}}))
	srv = loaders.Middleware(db, srv)

	http.Handle("/", playground.Handler("SaaS Sample API", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
