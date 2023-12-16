package e2etest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/99designs/gqlgen/client"
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

	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{}}))

	c := client.New(srv)

	c.MustPost(`mutation createTodo($content: String!) {createTodo(input: {content: $content}) {id}}`, &struct {
		CreateTodo *gql.Todo
	}{},
		client.Var("content", "test"),
	)

	var response struct {
		Todos []*gql.Todo
	}
	c.MustPost(`query { todos {id content done}}`, &response)

	if len(response.Todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(response.Todos))
	}

	if response.Todos[0].Content != "test" {
		t.Fatalf("expected todo content to be 'test', got %s", response.Todos[0].Content)
	}

	if response.Todos[0].Done {
		t.Fatalf("expected todo to be not done, got %v", response.Todos[0].Done)
	}

	c.MustPost(`mutation ToggleTodo($id: ID!, $content: String!) { updateTodo (input: {id: $id, content: $content}){id}}`, &struct {
		UpdateTodo *gql.Todo
	}{}, client.Var("id", response.Todos[0].ID), client.Var("content", "updated"))

	c.MustPost(`query { todos {id content done}}`, &response)

	if response.Todos[0].Content != "updated" {
		t.Fatalf("expected todo content to be 'updated', got %s", response.Todos[0].Content)
	}

	c.MustPost(`mutation ToggleTodo($id: ID!) { updateTodo (input: {id: $id, done: true}){id}}`, &struct {
		UpdateTodo *gql.Todo
	}{}, client.Var("id", response.Todos[0].ID))

	c.MustPost(`query { todos {id content done}}`, &response)

	if !response.Todos[0].Done {
		t.Fatalf("expected todo to be done, got %v", response.Todos[0].Done)
	}

	c.MustPost(`mutation DeleteTodo ($id: ID!) { deleteTodo (id: $id) }`, &struct {
		DeleteTodo string
	}{}, client.Var("id", response.Todos[0].ID))

	c.MustPost(`query { todos {id content done}}`, &response)

	if len(response.Todos) == 0 {
		t.Fatalf("expected 0 todo, got %d", len(response.Todos))
	}
}
