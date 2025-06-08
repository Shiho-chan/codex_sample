package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go backend!")
	})

	http.HandleFunc("/api/comments/", commentByIDHandler)
	http.HandleFunc("/api/comments", commentsHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
