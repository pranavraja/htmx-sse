package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	//go:embed templates/*.html
	templateFS embed.FS

	templates = template.Must(template.New("").ParseFS(templateFS, "templates/*.html"))
)

type quizHandler struct {
	question int64
	closed   *atomic.Bool

	sse  sseHandler
	auth authenticator
}

func (q *quizHandler) check(questionNumber int64, answer string) bool {
	switch questionNumber {
	case 1:
		return strings.EqualFold(answer, "adventure")
	case 2:
		return strings.Contains(strings.ToLower(answer), "zune")
	case 3:
		return strings.Contains(strings.ToLower(answer), "envelope")
	case 4:
		return strings.Contains(strings.ToLower(answer), "coin")
	case 5:
		return strings.EqualFold(answer, "c5")
	default:
		return false
	}
}

func (q *quizHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	questionNumber := atomic.LoadInt64(&q.question)
	switch r.Method {
	case http.MethodPost:
		answer := strings.TrimSpace(r.FormValue("answer"))
		name, err := q.auth.Username(r)
		if err != nil {
			http.Error(w, "couldn't authenticate you: "+err.Error(), http.StatusBadRequest)
			return
		}
		if name == "" {
			name = "anonymous"
		}
		log.Printf("question %d: attempt from %s: %s", questionNumber, name, answer)
		if q.closed.Load() {
			return
		}
		if q.check(questionNumber, answer) {
			q.closed.Store(true)
			q.sse.Broadcast(quizEvent{Type: "correct", From: name, Data: answer})
		} else {
			q.sse.Broadcast(quizEvent{Type: "wrong", From: name, Data: answer})
		}
	default:
		var quiz struct {
			Question int64
		}
		quiz.Question = questionNumber
		if err := templates.ExecuteTemplate(w, "questions.html", quiz); err != nil {
			log.Printf("failed to execute template: %s", err)
		}
	}
}

func (q *quizHandler) Admin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		q.closed.Store(false)
		atomic.AddInt64(&q.question, 1)
		q.sse.Broadcast(quizEvent{Type: "next"})
	case http.MethodPut:
		q.sse.Broadcast(quizEvent{Type: "reveal", Data: strconv.FormatInt(atomic.LoadInt64(&q.question), 10)})
	}
}
