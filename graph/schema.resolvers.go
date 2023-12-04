package graph

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/taro-28/saas-sample-api/db"
	"github.com/taro-28/saas-sample-api/graph/model"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
	todo := &model.Todo{
		Text: input.Text,
		ID:   randNumber.String(),
		User: &model.User{ID: input.UserID, Name: "user " + input.UserID},
	}

	db := db.Get()
	if _, err := db.Exec(
		"insert into todos (id, text) values (?, ?);",
		randNumber.Int64(),
		todo.Text,
	); err != nil {
		log.Fatalf("failed to insert: %v", err)
	}
	defer db.Close()

	return todo, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	db := db.Get()
	rows, err := db.Query("select id, text from todos;")
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	r.todos = []*model.Todo{}

	for rows.Next() {
		var id string
		var text string
		if err := rows.Scan(&id, &text); err != nil {
			log.Fatalf("failed to scan: %v", err)
		}
		fmt.Printf("id: %s, text: %s\n", id, text)

		r.todos = append(r.todos, &model.Todo{
			ID:   id,
			Text: text,
		})
	}

	defer rows.Close()
	defer db.Close()

	return r.todos, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
