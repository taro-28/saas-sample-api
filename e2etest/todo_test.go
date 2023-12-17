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
	mysqlContainer, err := mysql.RunContainer(ctx,
		mysql.WithDatabase("foo"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	if err != nil {
		panic(err)
	}

	ep, err := mysqlContainer.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}

	originalValue := os.Getenv("DSN")
	os.Setenv("DSN", fmt.Sprintf("root:password@tcp(%s)/foo", ep))

	t.Cleanup(func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
		os.Setenv("DSN", originalValue)
	})

	sqlFileContent, err := os.ReadFile("../db/schema.sql")
	if err != nil {
		t.Fatal(err)
	}

	db.Get().Exec(string(sqlFileContent))
}

func setupGqlServerAndClient(t *testing.T) *Client {
	t.Helper()
	h := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{}}))
	s := httptest.NewServer(h)

	t.Cleanup(func() {
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

	want := &TodoTest{
		Todos: []*TodoFragment{
			{
				ID:        createRes.CreateTodo.ID,
				Content:   "test",
				Done:      false,
				CreatedAt: int(testtime.Now().Unix()),
			},
		},
	}

	if createRes.CreateTodo.ID == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err := gqlClient.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, todosRes); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	updateContentRes, err := gqlClient.UpdateTodoContent(ctx, todosRes.Todos[0].ID, "updated")
	if err != nil {
		t.Fatal(err)
	}

	if updateContentRes.UpdateTodo.ID == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err = gqlClient.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if todosRes.Todos[0].Content != "updated" {
		t.Fatalf("expected todo content to be 'updated', got %s", todosRes.Todos[0].Content)
	}

	completeRes, err := gqlClient.CompleteTodo(ctx, todosRes.Todos[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if completeRes.UpdateTodo.ID == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err = gqlClient.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !todosRes.Todos[0].Done {
		t.Fatalf("expected todo to be done, got %v", todosRes.Todos[0].Done)
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

	if len(todosRes.Todos) != 0 {
		t.Fatalf("expected 0 todo, got %d", len(todosRes.Todos))
	}
}
