package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.41

import (
	"context"
	"log"
	"time"

	"github.com/rs/xid"
	"github.com/taro-28/saas-sample-api/db"
	"github.com/taro-28/saas-sample-api/models"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input CreateTodoInput) (*Todo, error) {
	todo := &models.Todo{
		ID:        xid.New().String(),
		Content:   input.Content,
		Done:      false,
		CreatedAt: uint(time.Now().Unix()),
	}

	db := db.Get()
	if err := todo.Insert(ctx, db); err != nil {
		log.Fatalf("failed to insert: %v", err)
	}
	defer db.Close()

	return &Todo{
		ID:        todo.ID,
		Content:   todo.Content,
		Done:      todo.Done,
		CreatedAt: int(todo.CreatedAt),
	}, nil
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, input UpdateTodoInput) (*Todo, error) {
	db := db.Get()

	todo, err := models.TodoByID(ctx, db, input.ID)
	if err != nil {
		log.Fatalf("failed to get todo by id: %v", err)
	}

	if input.Content != nil {
		todo.Content = *input.Content
	}
	if input.Done != nil {
		todo.Done = *input.Done
	}

	if err := todo.Update(ctx, db); err != nil {
		log.Fatalf("failed to update todo: %v", err)
	}

	defer db.Close()

	return &Todo{
		ID:        todo.ID,
		Content:   todo.Content,
		Done:      todo.Done,
		CreatedAt: int(todo.CreatedAt),
	}, nil
}

// DeleteTodo is the resolver for the deleteTodo field.
func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (string, error) {
	db := db.Get()

	todo, err := models.TodoByID(ctx, db, id)
	if err != nil {
		log.Fatalf("failed to get todo by id: %v", err)
	}

	err = todo.Delete(ctx, db)
	if err != nil {
		log.Fatalf("failed to delete todo: %v", err)
	}

	defer db.Close()

	return id, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*Todo, error) {
	db := db.Get()

	todos, err := models.AllTodos(ctx, db)
	if err != nil {
		log.Fatalf("failed to get all todos: %v", err)
	}

	defer db.Close()

	var gqlTodos []*Todo
	for _, todo := range todos {
		gqlTodos = append(gqlTodos, &Todo{
			ID:        todo.ID,
			Content:   todo.Content,
			Done:      todo.Done,
			CreatedAt: int(todo.CreatedAt),
		})
	}

	return gqlTodos, nil
}
