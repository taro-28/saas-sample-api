package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.41

import (
	"context"
	"time"

	"github.com/morikuni/failure"
	"github.com/rs/xid"
	"github.com/taro-28/saas-sample-api/fail"
	gql "github.com/taro-28/saas-sample-api/gql/model"
	"github.com/taro-28/saas-sample-api/models"
)

// Todos is the resolver for the todos field.
func (r *categoryResolver) Todos(ctx context.Context, obj *gql.Category) ([]*gql.Todo, error) {
	todo, err := models.AllTodos(ctx, r.DB)
	if err != nil {
		return nil, failure.Translate(err, fail.InternalServerError, failure.Message("failed to get todos"))
	}

	var todos []*gql.Todo
	for _, t := range todo {
		if t.CategoryID.String != obj.ID {
			continue
		}
		todos = append(todos, &gql.Todo{
			ID:        t.ID,
			Content:   t.Content,
			Done:      t.Done,
			CreatedAt: int(t.CreatedAt),
		})
	}

	return todos, nil
}

// CreateCategory is the resolver for the createCategory field.
func (r *mutationResolver) CreateCategory(ctx context.Context, input gql.CreateCategoryInput) (*gql.Category, error) {
	category := &models.Category{
		ID:        xid.New().String(),
		Name:      input.Name,
		CreatedAt: uint(time.Now().Unix()),
	}

	if err := category.Insert(ctx, r.DB); err != nil {
		return nil, failure.Translate(err, fail.InternalServerError, failure.Message("failed to insert category"))
	}

	return &gql.Category{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: int(category.CreatedAt),
	}, nil
}

// UpdateCategory is the resolver for the updateCategory field.
func (r *mutationResolver) UpdateCategory(ctx context.Context, input gql.UpdateCategoryInput) (*gql.Category, error) {
	category, err := models.CategoryByID(ctx, r.DB, input.ID)
	if err != nil {
		return nil, failure.Translate(err, fail.NotFound, failure.Messagef("category not found by id: %s", input.ID))
	}

	category.Name = input.Name
	if err := category.Update(ctx, r.DB); err != nil {
		return nil, failure.Translate(err, fail.InternalServerError, failure.Message("failed to update category"))
	}

	return &gql.Category{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: int(category.CreatedAt),
	}, nil
}

// DeleteCategory is the resolver for the deleteCategory field.
func (r *mutationResolver) DeleteCategory(ctx context.Context, id string) (string, error) {
	category, err := models.CategoryByID(ctx, r.DB, id)
	if err != nil {
		return "", failure.Translate(err, fail.NotFound, failure.Messagef("category not found by id: %s", id))
	}

	if err = category.Delete(ctx, r.DB); err != nil {
		return "", failure.Translate(err, fail.InternalServerError, failure.Message("failed to delete category"))
	}

	return id, nil
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context) ([]*gql.Category, error) {
	categories, err := models.AllCategorys(ctx, r.DB)
	if err != nil {
		return nil, failure.Translate(err, fail.InternalServerError, failure.Message("failed to get categories"))
	}

	var gqlCategories []*gql.Category
	for _, c := range categories {
		gqlCategories = append(gqlCategories, &gql.Category{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: int(c.CreatedAt),
		})
	}
	return gqlCategories, nil
}

// Category returns CategoryResolver implementation.
func (r *Resolver) Category() CategoryResolver { return &categoryResolver{r} }

type categoryResolver struct{ *Resolver }
