package models

type Edge[T any] struct {
	Node   *T      `json:"node"`
	Cursor *string `json:"cursor"`
}

type Connection[T any] struct {
	Edges    []*Edge[T] `json:"edges"`
	PageInfo *PageInfo  `json:"pageInfo"`
}
