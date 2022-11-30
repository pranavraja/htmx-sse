package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

var (
	//go:embed templates/*.html
	templateFS embed.FS

	templates = template.Must(template.New("").ParseFS(templateFS, "templates/*.html"))
)

type quizHandler struct {
	question int64

	sse  sseHandler
	auth authenticator
}

func (q *quizHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		answer := r.FormValue("answer")
		name, err := q.auth.Username(r)
		if err != nil {
			http.Error(w, "couldn't authenticate you: "+err.Error(), http.StatusBadRequest)
			return
		}
		if name == "" {
			name = "anonymous"
		}
		q.sse.Broadcast(quizEvent{Type: "attempt", From: name, Data: answer})
	default:
		var quiz struct {
			Question int64
		}
		quiz.Question = atomic.LoadInt64(&q.question)
		if err := templates.ExecuteTemplate(w, "questions.html", quiz); err != nil {
			log.Printf("failed to execute template: %s", err)
		}
	}
}

func (q *quizHandler) Admin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		atomic.AddInt64(&q.question, 1)
		q.sse.Broadcast(quizEvent{Type: "next"})
	case http.MethodPut:
		q.sse.Broadcast(quizEvent{Type: "reveal", Data: strconv.FormatInt(atomic.LoadInt64(&q.question), 10)})
	}
}
