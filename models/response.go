package models

type Response[T BaseModel] struct {
	Content []T `json:"content"`
}
