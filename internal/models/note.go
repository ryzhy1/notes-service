package models

type Note struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Owner   string `json:"owner"`
}
