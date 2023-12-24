package gql

type Todo struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	Done       bool   `json:"done"`
	CreatedAt  int    `json:"createdAt"`
	CategoryID string `json:"-"`
}
