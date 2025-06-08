package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Comment struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var (
	comments   []Comment
	nextID     int
	commentsMu sync.Mutex
)

func init() {
	nextID = 1
}

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
		id := nextID
		nextID++
		comments = append(comments, Comment{ID: id, Text: c.Text})
		commentsMu.Unlock()
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func commentByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	commentsMu.Lock()
	defer commentsMu.Unlock()
	index := -1
	for i, c := range comments {
		if c.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodDelete:
		comments = append(comments[:index], comments[index+1:]...)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
