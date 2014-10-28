package main

import (
	"html/template"
	"net/http"
	"path"
)

type Task struct {
	Name  string
	Description string
}
func main() {
	http.HandleFunc("/", ShowTasks)
	http.ListenAndServe(":8080", nil)
}
func ShowTasks(w http.ResponseWriter, r *http.Request) {

	tasks:=[]Task {
		{"Task1", "Task Desc"},
		{ "Task2","Task Desc"},
	}
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
