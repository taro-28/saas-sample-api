package e2etest

// import (
// 	"context"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// 	gql "github.com/taro-28/saas-sample-api/gql/model"
// 	"github.com/tenntenn/testtime"
// )

// func TestE2E_Category(t *testing.T) {
// 	t.Parallel()
// 	ctx := context.Background()
// 	setupDB(ctx, t)
// 	gqlClient := setupGqlServerAndClient(t)

// 	createRes, err := gqlClient.CreateCategory(ctx, "test")
// 	if err != nil {
// 		t.Fatalf("failed to create category: %v", err)
// 	}

// 	wantCreated := &CategoryFragment{
// 		ID:        createRes.CreateCategory.ID,
// 		Name:      "test",
// 		CreatedAt: int(testtime.Now().Unix()),
// 		Todos:     []*TodoFragment{},
// 	}

// 	if diff := cmp.Diff(wantCreated, &createRes.CreateCategory); diff != "" {
// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
// 	}

// 	categoriesRes, err := gqlClient.Categories(ctx)
// 	if err != nil {
// 		t.Fatalf("failed to get categories: %v", err)
// 	}

// 	wantList := &Categories{Categories: []*CategoryFragment{wantCreated}}
// 	if diff := cmp.Diff(wantList, categoriesRes); diff != "" {
// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
// 	}

// 	updateRes, err := gqlClient.UpdateCategory(ctx, gql.UpdateCategoryInput{
// 		ID:   wantCreated.ID,
// 		Name: "updated",
// 	})
// 	if err != nil {
// 		t.Fatalf("failed to update category name: %v", err)
// 	}

// 	wantUpdated := &CategoryFragment{
// 		ID:        updateRes.UpdateCategory.ID,
// 		Name:      "updated",
// 		CreatedAt: int(testtime.Now().Unix()),
// 		Todos:     []*TodoFragment{},
// 	}

// 	if diff := cmp.Diff(wantUpdated, &updateRes.UpdateCategory); diff != "" {
// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
// 	}

// 	categoriesRes, err = gqlClient.Categories(ctx)
// 	if err != nil {
// 		t.Fatalf("failed to get categories: %v", err)
// 	}

// 	wantList = &Categories{Categories: []*CategoryFragment{wantUpdated}}
// 	if diff := cmp.Diff(wantList, categoriesRes); diff != "" {
// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
// 	}

// 	deleteRes, err := gqlClient.DeleteCategory(ctx, wantUpdated.ID)
// 	if err != nil {
// 		t.Fatalf("failed to delete category: %v", err)
// 	}
// 	if deleteRes.DeleteCategory == "" {
// 		t.Fatalf("failed to delete category: %v", err)
// 	}

// 	categoriesRes, err = gqlClient.Categories(ctx)
// 	if err != nil {
// 		t.Fatalf("failed to get categories: %v", err)
// 	}

// 	wantList = &Categories{Categories: []*CategoryFragment{}}
// 	if diff := cmp.Diff(wantList, categoriesRes); diff != "" {
// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
// 	}
// }
