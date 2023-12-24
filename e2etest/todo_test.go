package e2etest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/taro-28/saas-sample-api/db"
	"github.com/taro-28/saas-sample-api/gql"
	"github.com/tenntenn/testtime"
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

func TestE2E_Todo(t *testing.T) {
	ctx := context.Background()

	setupDB(ctx, t)
	gqlClient := setupGqlServerAndClient(t)

	createRes, err := gqlClient.CreateTodo(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	wantCreated := &TodoFragment{
		ID:        createRes.CreateTodo.ID,
		Content:   "test",
		Done:      false,
		CreatedAt: int(testtime.Now().Unix()),
	}
	if diff := cmp.Diff(wantCreated, &createRes.CreateTodo); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	todosRes, err := gqlClient.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wantList := &TodoTest{Todos: []*TodoFragment{wantCreated}}
	if diff := cmp.Diff(wantList, todosRes); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	updateContentRes, err := gqlClient.UpdateTodoContent(ctx, wantCreated.ID, "updated")
	if err != nil {
		t.Fatal(err)
	}

	wantUpdated := &TodoFragment{
		ID:        updateContentRes.UpdateTodo.ID,
		Content:   "updated",
		Done:      false,
		CreatedAt: int(testtime.Now().Unix()),
	}
	if diff := cmp.Diff(wantUpdated, &updateContentRes.UpdateTodo); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	todosRes, err = gqlClient.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wantList = &TodoTest{Todos: []*TodoFragment{wantUpdated}}
	if diff := cmp.Diff(wantList, todosRes); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	completeRes, err := gqlClient.CompleteTodo(ctx, todosRes.Todos[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	wantCompleted := &TodoFragment{
		ID:        completeRes.UpdateTodo.ID,
		Content:   "updated",
		Done:      true,
		CreatedAt: int(testtime.Now().Unix()),
	}
	if diff := cmp.Diff(wantCompleted, &completeRes.UpdateTodo); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	todosRes, err = gqlClient.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wantList = &TodoTest{Todos: []*TodoFragment{wantCompleted}}
	if diff := cmp.Diff(wantList, todosRes); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	deleteRes, err := gqlClient.DeleteTodo(ctx, todosRes.Todos[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if deleteRes.DeleteTodo == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err = gqlClient.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wantList = &TodoTest{Todos: []*TodoFragment{}}
	if diff := cmp.Diff(wantList, todosRes); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
}
