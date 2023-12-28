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

	t.Run("delete", func(t *testing.T) {
		t.Parallel()
		createRes, err := gqlClient.CreateCategory(ctx, gql.CreateCategoryInput{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("failed to create category: %v", err)
		}
		id := createRes.CreateCategory.ID

		testCases := map[string]struct {
			input   string
			want    string
			wantErr bool
		}{
			"ok:basic": {
				input: id,
				want:  id,
			},
			"ng: invalid id": {
				input:   "invalid",
				wantErr: true,
			},
		}

		for name, tc := range testCases {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				deleteRes, err := gqlClient.DeleteCategory(ctx, tc.input)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("want error but got nil")
					}

					return
				}

				if err != nil {
					t.Fatalf("failed to delete category: %v", err)
				}

				if deleteRes.DeleteCategory != tc.want {
					t.Fatalf("want %s but got %s", tc.want, deleteRes.DeleteCategory)
				}

				categoriesRes, err := gqlClient.Categories(ctx)
				if err != nil {
					t.Fatalf("failed to get categories: %v", err)
				}

				for _, category := range categoriesRes.Categories {
					if category.ID != createRes.CreateCategory.ID {
						continue
					}

					if diff := cmp.Diff(tc.want, category.ID); diff != "" {
						t.Fatalf("mismatch (-want +got):\n%s", diff)
					}
				}
			})
		}
	})
}
