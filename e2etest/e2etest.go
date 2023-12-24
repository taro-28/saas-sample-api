package e2etest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/taro-28/saas-sample-api/db"
	"github.com/taro-28/saas-sample-api/gql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func setupDB(ctx context.Context, t *testing.T) {
	t.Helper()
	const (
		DBName = "foo"
		DBUser = "root"
		DBPass = "password"
	)

	mysqlContainer, err := mysql.RunContainer(ctx,
		mysql.WithDatabase(DBName),
		mysql.WithUsername(DBUser),
		mysql.WithPassword(DBPass),
		mysql.WithScripts("../db/schema.sql"),
	)
	if err != nil {
		panic(err)
	}

	ep, err := mysqlContainer.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}

	originalValue := os.Getenv("DSN")
	os.Setenv("DSN", fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, ep, DBName))

	t.Cleanup(func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
		os.Setenv("DSN", originalValue)
	})
}

func setupGqlServerAndClient(t *testing.T) *Client {
	t.Helper()

	db, cleanup, err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	h := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{
		DB: db,
	}}))
	s := httptest.NewServer(h)

	t.Cleanup(func() {
		cleanup()
		s.Close()
	})

	return NewClient(http.DefaultClient, s.URL)
}
