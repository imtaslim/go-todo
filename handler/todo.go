package handler

import (
	"net/http"
	"strconv"
)

type FormData struct {
	ID int
	Task string
	Title string
	Errors map[string]string
}

func (h *Handler) TodoCreate (rw http.ResponseWriter, r *http.Request) {
	vErrs := map[string]string{"task": "", "title": ""}
	task := ""
	title := ""
	h.createFormData(rw, task, title, vErrs)
	return
}

func (h *Handler) TodoStore (rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	task := r.FormValue("Task")
	title := r.FormValue("Title")
	if task == "" {
		rerrs := map[string]string{"task": "This field is required"}
		h.createFormData(rw, task, title, rerrs)
		return
	}
	if title == "" {
		rerrs := map[string]string{"title": "This field is required"}
		h.createFormData(rw, task, title, rerrs)
		return
	}
	if len(task) < 3 {
		rerrs := map[string]string{"task": "This field must be greater than or equals 3"}
		h.createFormData(rw, task, title, rerrs)
		return
	}
	if len(title) < 3 {
		rerrs := map[string]string{"title": "This field must be greater than or equals 3"}
		h.createFormData(rw, task, title, rerrs)
		return
	}

	const insertTodo = `INSERT INTO todos(title, task, is_completed) VALUES ($1, $2, $3)`
	res := h.db.MustExec(insertTodo, task, title, false)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// h.todos = append(h.todos, Todo{Task: task})
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) TodoComplete (rw http.ResponseWriter, r *http.Request) {
	Id := r.URL.Path[len("/todos/complete/"):]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const updateStatusTodo = `UPDATE todos SET is_completed = true WHERE id=$1`
	res := h.db.MustExec(updateStatusTodo, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) TodoEdit (rw http.ResponseWriter, r *http.Request) {
	Id := r.URL.Path[len("/todos/edit/"):]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getTodo = `SELECT * FROM todos WHERE id=$1`
	var todo Todo
	h.db.Get(&todo, getTodo, Id)
	
	rerrs := map[string]string{"title": "This field is required"}
		h.editFormData(rw, todo.ID, todo.Task, todo.Title, rerrs)
		return
}

func (h *Handler) TodoUpdate (rw http.ResponseWriter, r *http.Request) {
	Id := r.URL.Path[len("/todos/update/"):]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	t := r.FormValue("Task")
	ttl := r.FormValue("Title")
	id, err := strconv.Atoi(Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if t == "" {
		rerrs := map[string]string{"task": "This field is required"}
		h.editFormData(rw, id, t, ttl, rerrs)
		return
	}
	if ttl == "" {
		rerrs := map[string]string{"title": "This field is required"}
		h.editFormData(rw, id, t, ttl, rerrs)
		return
	}
	if len(t) < 3 {
		rerrs := map[string]string{"task": "This field must be greater than or equals 3"}
		h.editFormData(rw, id, t, ttl, rerrs)
		return
	}
	if len(ttl) < 3 {
		rerrs := map[string]string{"title": "This field must be greater than or equals 3"}
		h.editFormData(rw, id, t, ttl, rerrs)
		return
	}

	const updateStatusTodo = `UPDATE todos SET task=$1, title=$2 WHERE id=$3`
	res := h.db.MustExec(updateStatusTodo, t, ttl, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) TodoDelete (rw http.ResponseWriter, r *http.Request) {
	Id := r.URL.Path[len("/todos/delete/"):]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const deleteTodo = `DELETE FROM todos WHERE id=$1`
	res := h.db.MustExec(deleteTodo, Id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) createFormData (rw http.ResponseWriter, task string, title string, errs map[string]string) {
	form := FormData{
		Task: task,
		Title: title,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "create-todo.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) editFormData (rw http.ResponseWriter, id int, task string, title string, errs map[string]string) {
	form := FormData{
		ID: id,
		Task: task,
		Title: title,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "edit-todo.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}