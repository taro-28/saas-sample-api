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
	"github.com/taro-28/saas-sample-api/db"
	"github.com/taro-28/saas-sample-api/gql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestE2E_Todo(t *testing.T) {
	// mysqlのテストサーバーを起動する
	ctx := context.Background()

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

	// Clean up the container
	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()

	// 環境変数のDSNをテスト用のものに差し替えて元に戻す
	originalValue := os.Getenv("DSN")
	defer func() {
		os.Setenv("DSN", originalValue)
	}()
	os.Setenv("DSN", fmt.Sprintf("root:password@tcp(%s)/foo", ep))

	sqlFileContent, err := os.ReadFile("../db/schema.sql")
	if err != nil {
		t.Fatal(err)
	}

	// db/schema.sql を実行する
	db.Get().Exec(string(sqlFileContent))

	h := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{}}))
	s := httptest.NewServer(h)
	defer s.Close()

	c := NewClient(http.DefaultClient, s.URL)
	createRes, err := c.CreateTodo(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	if createRes.CreateTodo.ID == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err := c.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(todosRes.Todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todosRes.Todos))
	}

	if todosRes.Todos[0].Content != "test" {
		t.Fatalf("expected todo content to be 'test', got %s", todosRes.Todos[0].Content)
	}

	if todosRes.Todos[0].Done {
		t.Fatalf("expected todo to be not done, got %v", todosRes.Todos[0].Done)
	}

	updateContentRes, err := c.UpdateTodoContent(ctx, todosRes.Todos[0].ID, "updated")
	if err != nil {
		t.Fatal(err)
	}

	if updateContentRes.UpdateTodo.ID == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err = c.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if todosRes.Todos[0].Content != "updated" {
		t.Fatalf("expected todo content to be 'updated', got %s", todosRes.Todos[0].Content)
	}

	completeRes, err := c.CompleteTodo(ctx, todosRes.Todos[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if completeRes.UpdateTodo.ID == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err = c.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !todosRes.Todos[0].Done {
		t.Fatalf("expected todo to be done, got %v", todosRes.Todos[0].Done)
	}

	deleteRes, err := c.DeleteTodo(ctx, todosRes.Todos[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if deleteRes.DeleteTodo == "" {
		t.Fatal("expected todo id to be not empty")
	}

	todosRes, err = c.TodoTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(todosRes.Todos) != 0 {
		t.Fatalf("expected 0 todo, got %d", len(todosRes.Todos))
	}
}
