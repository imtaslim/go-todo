package main

import (
	//"database/sql"
	"log"
	"net/http"

	"todo/handler"

	_ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)



func main() {
	var schema = `
	CREATE TABLE IF NOT EXISTS todos (
		id serial,
		task text,
		title text,
		is_completed boolean,

		primary key(id)
	);
`
	db, err := sqlx.Connect("postgres", "user=postgres password=P@ssw0rd dbname=todo sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

	db.MustExec(schema)

	h := handler.New(db)

	http.HandleFunc("/", h.Home)

	http.HandleFunc("/todos/create", h.TodoCreate)
	http.HandleFunc("/todos/store", h.TodoStore)
	http.HandleFunc("/todos/complete/", h.TodoComplete)
	http.HandleFunc("/todos/edit/", h.TodoEdit)
	http.HandleFunc("/todos/update/", h.TodoUpdate)
	http.HandleFunc("/todos/delete/", h.TodoDelete)
	
	log.Println("Server Starting....")
	if err := http.ListenAndServe("127.0.0.1:3000", nil); err != nil {
		log.Fatal(err)
	}
}
