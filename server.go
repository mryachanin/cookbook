package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	data := "Hello World!"
	t := template.Must(template.ParseFiles(path.Join("templates", "index.html")))
	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
