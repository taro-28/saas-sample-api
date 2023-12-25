package loaders

// import graph gophers with your other imports
import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/graph-gophers/dataloader"
	gql "github.com/taro-28/saas-sample-api/gql/model"
	"github.com/taro-28/saas-sample-api/models"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// categoryReader reads Categorys from a database
type categoryReader struct {
	db *sql.DB
}

// getCategorys implements a batch function that can retrieve many categories by ID,
// for use in a dataloader
func (c *categoryReader) getCategories(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	categories, err := models.AllCategorys(ctx, c.db)
	if err != nil {
		return handleError[*gql.Category](len(keys), err)
	}

	// create a result for each key
	result := make([]*dataloader.Result, len(keys))
	for i, key := range keys {
		for _, category := range categories {
			if category.ID == key.String() {
				result[i] = &dataloader.Result{Data: &gql.Category{
					ID:        category.ID,
					Name:      category.Name,
					CreatedAt: int(category.CreatedAt),
				}}
				break
			}
		}
	}
	return result

}

// handleError creates array of result with the same error repeated for as many items requested
func handleError[T any](itemsLength int, err error) []*dataloader.Result {
	result := make([]*dataloader.Result, itemsLength)
	for i := 0; i < itemsLength; i++ {
		result[i] = &dataloader.Result{Error: err}
	}
	return result
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	CategoryLoader *dataloader.Loader
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *sql.DB) *Loaders {
	// define the data loader
	ur := &categoryReader{db: conn}
	return &Loaders{
		CategoryLoader: dataloader.NewBatchedLoader(ur.getCategories, dataloader.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(conn *sql.DB, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(conn)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetCategory returns single category by id efficiently
func GetCategory(ctx context.Context, id string) (*gql.Category, error) {
	loaders := For(ctx)
	result, err := loaders.CategoryLoader.Load(ctx, dataloader.StringKey(id))()
	return result.(*gql.Category), err
}

// GetCategories returns many categories by ids efficiently
func GetCategories(ctx context.Context, ids dataloader.Keys) ([]*gql.Category, []error) {
	loaders := For(ctx)
	result, err := loaders.CategoryLoader.LoadMany(ctx, ids)()
	if err != nil {
		return nil, err
	}

	var categories []*gql.Category
	for _, r := range result {
		categories = append(categories, r.(*gql.Category))
	}
	return categories, nil
}
