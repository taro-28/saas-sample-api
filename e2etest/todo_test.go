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
		category, err := gqlClient.CreateCategory(ctx, "test")
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

	// testCases := map[string]struct {
	// 	createInput gql.CreateTodoInput
	// 	wantCreated *TodoFragment
	// 	updateInput gql.UpdateTodoInput
	// 	wantUpdated *TodoFragment
	// }{
	// 	"basic": {
	// 		createInput: gql.CreateTodoInput{
	// 			Content:    "test",
	// 			CategoryID: &category.ID,
	// 		},
	// 		wantCreated: &TodoFragment{
	// 			ID:        "",
	// 			Content:   "test",
	// 			Done:      false,
	// 			CreatedAt: int(testtime.Now().Unix()),
	// 			Category: &struct {
	// 				ID        string "json:\"id\" graphql:\"id\""
	// 				Name      string "json:\"name\" graphql:\"name\""
	// 				CreatedAt int    "json:\"createdAt\" graphql:\"createdAt\""
	// 			}{
	// 				ID:        category.ID,
	// 				Name:      category.Name,
	// 				CreatedAt: category.CreatedAt,
	// 			},
	// 		},
	// 		updateInput: gql.UpdateTodoInput{
	// 			Content: "updated",
	// 		},
	// 		wantUpdated: &TodoFragment{
	// 			ID:        "",
	// 			Content:   "updated",
	// 			Done:      false,
	// 			CreatedAt: int(testtime.Now().Unix()),
	// 			Category: &struct {
	// 				ID        string "json:\"id\" graphql:\"id\""
	// 				Name      string "json:\"name\" graphql:\"name\""
	// 				CreatedAt int    "json:\"createdAt\" graphql:\"createdAt\""
	// 			}{
	// 				ID:        category.ID,
	// 				Name:      category.Name,
	// 				CreatedAt: category.CreatedAt,
	// 			},
	// 		},
	// 	},
	// 	"non category": {},
	// }

	// for name, tc := range testCases {
	// 	createRes, err := gqlClient.CreateTodo(ctx, "test")
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	wantCreated := &TodoFragment{
	// 		ID:        createRes.CreateTodo.ID,
	// 		Content:   "test",
	// 		Done:      false,
	// 		CreatedAt: int(testtime.Now().Unix()),
	// 	}
	// 	if diff := cmp.Diff(wantCreated, &createRes.CreateTodo); diff != "" {
	// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
	// 	}

	// 	todosRes, err := gqlClient.Todos(ctx)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	wantList := &Todos{Todos: []*TodoFragment{wantCreated}}
	// 	if diff := cmp.Diff(wantList, todosRes); diff != "" {
	// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
	// 	}

	// 	updateContentRes, err := gqlClient.UpdateTodoContent(ctx, wantCreated.ID, "updated")
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	wantUpdated := &TodoFragment{
	// 		ID:        updateContentRes.UpdateTodo.ID,
	// 		Content:   "updated",
	// 		Done:      false,
	// 		CreatedAt: int(testtime.Now().Unix()),
	// 	}
	// 	if diff := cmp.Diff(wantUpdated, &updateContentRes.UpdateTodo); diff != "" {
	// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
	// 	}

	// 	todosRes, err = gqlClient.Todos(ctx)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	wantList = &Todos{Todos: []*TodoFragment{wantUpdated}}
	// 	if diff := cmp.Diff(wantList, todosRes); diff != "" {
	// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
	// 	}

	// 	completeRes, err := gqlClient.CompleteTodo(ctx, todosRes.Todos[0].ID)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	wantCompleted := &TodoFragment{
	// 		ID:        completeRes.UpdateTodo.ID,
	// 		Content:   "updated",
	// 		Done:      true,
	// 		CreatedAt: int(testtime.Now().Unix()),
	// 	}
	// 	if diff := cmp.Diff(wantCompleted, &completeRes.UpdateTodo); diff != "" {
	// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
	// 	}

	// 	todosRes, err = gqlClient.Todos(ctx)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	wantList = &Todos{Todos: []*TodoFragment{wantCompleted}}
	// 	if diff := cmp.Diff(wantList, todosRes); diff != "" {
	// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
	// 	}

	// 	deleteRes, err := gqlClient.DeleteTodo(ctx, todosRes.Todos[0].ID)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	if deleteRes.DeleteTodo == "" {
	// 		t.Fatal("expected todo id to be not empty")
	// 	}

	// 	todosRes, err = gqlClient.Todos(ctx)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	wantList = &Todos{Todos: []*TodoFragment{}}
	// 	if diff := cmp.Diff(wantList, todosRes); diff != "" {
	// 		t.Fatalf("mismatch (-want +got):\n%s", diff)
	// 	}
	// }
}
