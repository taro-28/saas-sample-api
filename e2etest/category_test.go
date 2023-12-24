package e2etest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gql "github.com/taro-28/saas-sample-api/gql/model"
	"github.com/tenntenn/testtime"
)

func TestE2E_Category(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	setupDB(ctx, t)
	gqlClient := setupGqlServerAndClient(t)

	t.Run("create", func(t *testing.T) {
		t.Parallel()
		testCases := map[string]struct {
			input   gql.CreateCategoryInput
			want    *CategoryFragment
			wantErr bool
		}{
			"ok:basic": {
				input: gql.CreateCategoryInput{
					Name: "test",
				},
				want: &CategoryFragment{
					Name:      "test",
					CreatedAt: int(testtime.Now().Unix()),
				},
			},
		}

		for name, tc := range testCases {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				createRes, err := gqlClient.CreateCategory(ctx, tc.input)
				if err != nil {
					t.Fatalf("failed to create category: %v", err)
				}

				if diff := cmp.Diff(tc.want, &createRes.CreateCategory,
					cmpopts.IgnoreFields(CategoryFragment{}, "ID"),
				); diff != "" {
					t.Fatalf("mismatch (-want +got):\n%s", diff)
				}

				categoriesRes, err := gqlClient.Categories(ctx)
				if err != nil {
					t.Fatalf("failed to get categories: %v", err)
				}

				for _, category := range categoriesRes.Categories {
					if category.ID != createRes.CreateCategory.ID {
						continue
					}

					if diff := cmp.Diff(tc.want, category,
						cmpopts.IgnoreFields(CategoryFragment{}, "ID"),
					); diff != "" {
						t.Fatalf("mismatch (-want +got):\n%s", diff)
					}
				}
			})
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()
		testCases := map[string]struct {
			input   gql.UpdateCategoryInput
			want    *CategoryFragment
			wantErr bool
		}{
			"ok:basic": {
				input: gql.UpdateCategoryInput{
					Name: "updated",
				},
				want: &CategoryFragment{
					Name:      "updated",
					CreatedAt: int(testtime.Now().Unix()),
				},
			},
			"ng: invalid id": {
				input: gql.UpdateCategoryInput{
					ID:   "invalid",
					Name: "updated",
				},
				wantErr: true,
			},
		}

		for name, tc := range testCases {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				createRes, err := gqlClient.CreateCategory(ctx, gql.CreateCategoryInput{
					Name: "test",
				})
				if err != nil {
					t.Fatalf("failed to create category: %v", err)
				}

				if tc.input.ID == "" {
					tc.input.ID = createRes.CreateCategory.ID
				}
				updateRes, err := gqlClient.UpdateCategory(ctx, tc.input)

				if tc.wantErr {
					if err == nil {
						t.Fatalf("want error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("failed to update category name: %v", err)
				}

				if diff := cmp.Diff(tc.want, &updateRes.UpdateCategory,
					cmpopts.IgnoreFields(CategoryFragment{}, "ID"),
				); diff != "" {
					t.Fatalf("mismatch (-want +got):\n%s", diff)
				}

				categoriesRes, err := gqlClient.Categories(ctx)
				if err != nil {
					t.Fatalf("failed to get categories: %v", err)
				}

				for _, category := range categoriesRes.Categories {
					if category.ID != createRes.CreateCategory.ID {
						continue
					}

					if diff := cmp.Diff(tc.want, category,
						cmpopts.IgnoreFields(CategoryFragment{}, "ID"),
					); diff != "" {
						t.Fatalf("mismatch (-want +got):\n%s", diff)
					}
				}
			})
		}
	})

	// createRes, err := gqlClient.CreateCategory(ctx, "test")
	// if err != nil {
	// 	t.Fatalf("failed to create category: %v", err)
	// }

	// wantCreated := &CategoryFragment{
	// 	ID:        createRes.CreateCategory.ID,
	// 	Name:      "test",
	// 	CreatedAt: int(testtime.Now().Unix()),
	// 	Todos:     []*TodoFragment{},
	// }

	// if diff := cmp.Diff(wantCreated, &createRes.CreateCategory); diff != "" {
	// 	t.Fatalf("mismatch (-want +got):\n%s", diff)
	// }

	// categoriesRes, err := gqlClient.Categories(ctx)
	// if err != nil {
	// 	t.Fatalf("failed to get categories: %v", err)
	// }

	// wantList := &Categories{Categories: []*CategoryFragment{wantCreated}}
	// if diff := cmp.Diff(wantList, categoriesRes); diff != "" {
	// 	t.Fatalf("mismatch (-want +got):\n%s", diff)
	// }

	// updateRes, err := gqlClient.UpdateCategory(ctx, gql.UpdateCategoryInput{
	// 	ID:   wantCreated.ID,
	// 	Name: "updated",
	// })
	// if err != nil {
	// 	t.Fatalf("failed to update category name: %v", err)
	// }

	// wantUpdated := &CategoryFragment{
	// 	ID:        updateRes.UpdateCategory.ID,
	// 	Name:      "updated",
	// 	CreatedAt: int(testtime.Now().Unix()),
	// 	Todos:     []*TodoFragment{},
	// }

	// if diff := cmp.Diff(wantUpdated, &updateRes.UpdateCategory); diff != "" {
	// 	t.Fatalf("mismatch (-want +got):\n%s", diff)
	// }

	// categoriesRes, err = gqlClient.Categories(ctx)
	// if err != nil {
	// 	t.Fatalf("failed to get categories: %v", err)
	// }

	// wantList = &Categories{Categories: []*CategoryFragment{wantUpdated}}
	// if diff := cmp.Diff(wantList, categoriesRes); diff != "" {
	// 	t.Fatalf("mismatch (-want +got):\n%s", diff)
	// }

	// deleteRes, err := gqlClient.DeleteCategory(ctx, wantUpdated.ID)
	// if err != nil {
	// 	t.Fatalf("failed to delete category: %v", err)
	// }
	// if deleteRes.DeleteCategory == "" {
	// 	t.Fatalf("failed to delete category: %v", err)
	// }

	// categoriesRes, err = gqlClient.Categories(ctx)
	// if err != nil {
	// 	t.Fatalf("failed to get categories: %v", err)
	// }

	// wantList = &Categories{Categories: []*CategoryFragment{}}
	// if diff := cmp.Diff(wantList, categoriesRes); diff != "" {
	// 	t.Fatalf("mismatch (-want +got):\n%s", diff)
	// }
}
