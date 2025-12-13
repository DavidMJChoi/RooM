package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/posts/", postHandler)

	fmt.Println("Server starting on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a blog page.")
}
