package handler

import (
	"html/template"

	"github.com/jmoiron/sqlx"
)

type Todo struct{
	ID int `db:"id"`
	Task string `db:"task"`
	Title string `db:"title"`
	IsCompleted bool `db:"is_completed"`
}

type Handler struct{
	templates *template.Template
	db *sqlx.DB
}

func New(db *sqlx.DB) *Handler {
	h := &Handler{
		db: db,
	}
	
	h.parseTemplates()

	return h
}

func (h *Handler) parseTemplates() {
	h.templates = template.Must(template.ParseFiles(
		"templates/create-todo.html",
		"templates/list-todo.html",
		"templates/edit-todo.html",
	))
}