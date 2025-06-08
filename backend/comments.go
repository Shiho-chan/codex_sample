package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Comment struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var (
	comments   []Comment
	commentsMu sync.Mutex
)

func commentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		commentsMu.Lock()
		defer commentsMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	case http.MethodPost:
		var c struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		commentsMu.Lock()
		id := len(comments) + 1
		comments = append(comments, Comment{ID: id, Text: c.Text})
		commentsMu.Unlock()
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
