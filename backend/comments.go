package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Comment struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

var (
	comments   []Comment
	nextID     int
	commentsMu sync.Mutex

	users = map[string]struct {
		Password string
		Role     string
	}{
		"admin": {Password: "password", Role: "admin"},
		"user":  {Password: "password", Role: "user"},
	}

	sessions   = make(map[string]string)
	sessionsMu sync.Mutex
)

func init() {
	nextID = 1
}

func newToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func getUser(r *http.Request) (string, string, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", "", false
	}
	sessionsMu.Lock()
	username, ok := sessions[cookie.Value]
	sessionsMu.Unlock()
	if !ok {
		return "", "", false
	}
	user, ok := users[username]
	if !ok {
		return "", "", false
	}
	return username, user.Role, true
}

func commentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		commentsMu.Lock()
		defer commentsMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	case http.MethodPost:
		username, _, ok := getUser(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var c struct {
			Name string `json:"name"`
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		commentsMu.Lock()
		id := nextID
		nextID++
		comments = append(comments, Comment{
			ID:        id,
			Username:  username,
			Name:      c.Name,
			Text:      c.Text,
			Timestamp: time.Now().Format(time.RFC3339),
		})
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
		username, role, ok := getUser(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if role != "admin" && comments[index].Username != username {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		comments = append(comments[:index], comments[index+1:]...)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	u, ok := users[creds.Username]
	if !ok || u.Password != creds.Password {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token := newToken()
	sessionsMu.Lock()
	sessions[token] = creds.Username
	sessionsMu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: "session", Value: token, Path: "/", HttpOnly: true})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}{creds.Username, u.Role})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	username, role, ok := getUser(r)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}{username, role})
}
