package e2etest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gql "github.com/taro-28/saas-sample-api/gql/model"
	"github.com/tenntenn/testtime"
)

func TestE2E_Todo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	setupDB(ctx, t)
	gqlClient := setupGqlServerAndClient(t)

	category := func() *CategoryFragment {
		category, err := gqlClient.CreateCategory(ctx, gql.CreateCategoryInput{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("failed to create category: %v", err)
		}
		return &category.CreateCategory
	}()

	t.Run("create", func(t *testing.T) {
		t.Parallel()
		testCases := map[string]struct {
			input   gql.CreateTodoInput
			want    *TodoFragment
			wantErr bool
		}{
			"ok:basic": {
				input: gql.CreateTodoInput{
					Content:    "test",
					CategoryID: &category.ID,
				},
				want: &TodoFragment{
					Content:   "test",
					Done:      false,
					CreatedAt: int(testtime.Now().Unix()),
					Category:  category,
				},
			},
			"ok:non category": {
				input: gql.CreateTodoInput{
					Content: "test",
				},
				want: &TodoFragment{
					Content:   "test",
					Done:      false,
					Category:  nil,
					CreatedAt: int(testtime.Now().Unix()),
				},
			},
			"ng:invalid category": {
				input: gql.CreateTodoInput{
					Content:    "test",
					CategoryID: func() *string { s := "invalid"; return &s }(),
				},
				wantErr: true,
			},
		}

		for name, tc := range testCases {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				createRes, err := gqlClient.CreateTodo(ctx, tc.input)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("failed to create todo: %v", err)
				}

				if diff := cmp.Diff(tc.want, &createRes.CreateTodo,
					cmpopts.IgnoreFields(TodoFragment{}, "ID"),
				); diff != "" {
					t.Fatalf("mismatch (-want +got):\n%s", diff)
				}

				todosRes, err := gqlClient.Todos(ctx)
				if err != nil {
					t.Fatalf("failed to get todos: %v", err)
				}

				for _, todo := range todosRes.Todos {
					if todo.ID != createRes.CreateTodo.ID {
						continue
					}

					if diff := cmp.Diff(tc.want, todo, cmpopts.IgnoreFields(TodoFragment{}, "ID")); diff != "" {
						t.Fatalf("mismatch (-want +got):\n%s", diff)
					}
				}
			})
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()
		testCases := map[string]struct {
			input   gql.UpdateTodoInput
			want    *TodoFragment
			wantErr bool
		}{
			"ok:basic": {
				input: gql.UpdateTodoInput{
					Content:    "updated",
					CategoryID: &category.ID,
				},
				want: &TodoFragment{
					Content:   "updated",
					Done:      false,
					CreatedAt: int(testtime.Now().Unix()),
					Category:  category,
				},
			},
			"ok:non category": {
				input: gql.UpdateTodoInput{
					Content: "updated",
				},
				want: &TodoFragment{
					Content:   "updated",
					Done:      false,
					CreatedAt: int(testtime.Now().Unix()),
				},
			},
			"ng:invalid id": {
				input: gql.UpdateTodoInput{
					ID:      "invalid",
					Content: "updated",
				},
				wantErr: true,
			},
		}

		for name, tc := range testCases {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				createRes, err := gqlClient.CreateTodo(ctx, gql.CreateTodoInput{
					Content:    "test",
					CategoryID: &category.ID,
				})
				if err != nil {
					t.Fatalf("failed to create todo: %v", err)
				}

				if tc.input.ID == "" {
					tc.input.ID = createRes.CreateTodo.ID
				}
				updateRes, err := gqlClient.UpdateTodo(ctx, tc.input)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("failed to update todo content: %v", err)
				}

				if diff := cmp.Diff(tc.want, &updateRes.UpdateTodo,
					cmpopts.IgnoreFields(TodoFragment{}, "ID"),
				); diff != "" {
					t.Fatalf("mismatch (-want +got):\n%s", diff)
				}

				todosRes, err := gqlClient.Todos(ctx)
				if err != nil {
					t.Fatalf("failed to get todos: %v", err)
				}

				for _, todo := range todosRes.Todos {
					if todo.ID != createRes.CreateTodo.ID {
						continue
					}

					if diff := cmp.Diff(tc.want, todo, cmpopts.IgnoreFields(TodoFragment{}, "ID")); diff != "" {
						t.Fatalf("mismatch (-want +got):\n%s", diff)
					}
				}
			})
		}

	})

	t.Run("update done", func(t *testing.T) {
		t.Parallel()
		testCases := map[string]struct {
			input   gql.UpdateTodoDoneInput
			want    *TodoFragment
			wantErr bool
		}{
			"ok: to done": {
				input: gql.UpdateTodoDoneInput{
					Done: true,
				},
				want: &TodoFragment{
					Content:   "test",
					Done:      true,
					CreatedAt: int(testtime.Now().Unix()),
					Category:  category,
				},
			},
			"ok: to not done": {
				input: gql.UpdateTodoDoneInput{
					Done: false,
				},
				want: &TodoFragment{
					Content:   "test",
					Done:      false,
					CreatedAt: int(testtime.Now().Unix()),
					Category:  category,
				},
			},
			"ng:invalid id": {
				input: gql.UpdateTodoDoneInput{
					ID:   "invalid",
					Done: true,
				},
				wantErr: true,
			},
		}

		for name, tc := range testCases {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				createRes, err := gqlClient.CreateTodo(ctx, gql.CreateTodoInput{
					Content:    "test",
					CategoryID: &category.ID,
				})
				if err != nil {
					t.Fatalf("failed to create todo: %v", err)
				}

				if tc.input.ID == "" {
					tc.input.ID = createRes.CreateTodo.ID
				}
				updateRes, err := gqlClient.UpdateTodoDone(ctx, tc.input)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("failed to update todo done: %v", err)
				}

				if diff := cmp.Diff(tc.want, &updateRes.UpdateTodoDone,
					cmpopts.IgnoreFields(TodoFragment{}, "ID"),
				); diff != "" {
					t.Fatalf("mismatch (-want +got):\n%s", diff)
				}

				todosRes, err := gqlClient.Todos(ctx)
				if err != nil {
					t.Fatalf("failed to get todos: %v", err)
				}

				for _, todo := range todosRes.Todos {
					if todo.ID != createRes.CreateTodo.ID {
						continue
					}

					if diff := cmp.Diff(tc.want, todo, cmpopts.IgnoreFields(TodoFragment{}, "ID")); diff != "" {
						t.Fatalf("mismatch (-want +got):\n%s", diff)
					}
				}
			})
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()
		createRes, err := gqlClient.CreateTodo(ctx, gql.CreateTodoInput{
			Content:    "test",
			CategoryID: &category.ID,
		})
		if err != nil {
			t.Fatalf("failed to create todo: %v", err)
		}

		id := createRes.CreateTodo.ID

		testCases := map[string]struct {
			input   string
			want    string
			wantErr bool
		}{
			"ok": {
				input: id,
				want:  id,
			},
			"ng:invalid id": {
				input:   "invalid",
				wantErr: true,
			},
		}

		for name, tc := range testCases {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				deleteRes, err := gqlClient.DeleteTodo(ctx, tc.input)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("failed to delete todo: %v", err)
				}

				if deleteRes.DeleteTodo != tc.want {
					t.Fatalf("want %s but got %s", tc.want, deleteRes.DeleteTodo)
				}

				todosRes, err := gqlClient.Todos(ctx)
				if err != nil {
					t.Fatalf("failed to get todos: %v", err)
				}

				for _, todo := range todosRes.Todos {
					if todo.ID != createRes.CreateTodo.ID {
						continue
					}

					t.Fatalf("want todo to be deleted but got %v", todo)
				}
			})
		}
	})
}
