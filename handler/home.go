package handler

import (
	"net/http"
)

type ListTodo struct{
	Todos []Todo
}

func (h *Handler) Home (rw http.ResponseWriter, r *http.Request) {
	todos := []Todo{}
    h.db.Select(&todos, "SELECT * FROM todos")
	lt := ListTodo{
		Todos: todos,
	}
	if err := h.templates.ExecuteTemplate(rw, "list-todo.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}