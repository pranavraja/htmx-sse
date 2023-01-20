package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

type quizHandler struct {
	question int64
	closed   *atomic.Bool

	sse  sseHandler
	auth authenticator
}

func (q *quizHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	questionNumber := atomic.LoadInt64(&q.question)
	name, err := q.auth.Username(r)
	if err != nil {
		if q.admin(r) {
			// we'll allow it since you're admin
			name = "admin"
		} else {
			http.Error(w, "couldn't authenticate you: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		answer := strings.Join(r.Form["answer"], ",")
		if len(answer) > 200 {
			answer = answer[:200]
		}
		log.Printf("question %d: guess from %s: %s", questionNumber, name, answer)
		if q.closed.Load() {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if questionNumber < 0 || questionNumber >= int64(len(questions)) {
			questionNumber = int64(len(questions)) - 1
		}
		if questions[questionNumber].check(answer) {
			q.closed.Store(true)
			event := quizEvent{QuestionNumber: questionNumber, Type: "correct", From: name, Data: answer}
			q.sse.Broadcast(event)
			if err := templates.ExecuteTemplate(w, "feedback-correct.html", event); err != nil {
				log.Printf("failed to execute template: %s", err)
			}
		} else {
			event := quizEvent{QuestionNumber: questionNumber, Type: "wrong", From: name, Data: answer}
			q.sse.Broadcast(event)
			if err := templates.ExecuteTemplate(w, "feedback-wrong.html", event); err != nil {
				log.Printf("failed to execute template: %s", err)
			}
		}
	default:
		var quiz struct {
			Question int64
			Admin    bool
		}
		quiz.Question = questionNumber
		quiz.Admin = q.admin(r)
		if err := templates.ExecuteTemplate(w, "quiz.html", quiz); err != nil {
			log.Printf("failed to execute template: %s", err)
		}
	}
}

func (q *quizHandler) admin(r *http.Request) bool {
	// Easy way to share the same UI but have admin permissions (e.g. to reveal answers) running locally.
	// Other people can get sent a public IP address (either via ngrok or in a local network),
	// so these will return `false`.
	isLocal := strings.HasPrefix(r.RemoteAddr, "127.0.0.1:") || strings.HasPrefix(r.RemoteAddr, "[::1]:")
	// Additional check if you have sent people URLs via ngrok, because ngrok is "technically" a local IP
	// We get around this by also checking that the Host header is localhost
	return isLocal && strings.HasPrefix(r.Host, "localhost:")
}

func (q *quizHandler) Admin(w http.ResponseWriter, r *http.Request) {
	if !q.admin(r) {
		http.Error(w, "sorry, you are not admin", http.StatusForbidden)
		return
	}
	switch r.Method {
	case http.MethodPost:
		q.closed.Store(false)
		atomic.AddInt64(&q.question, 1)
		q.sse.Broadcast(quizEvent{Type: "next"})
	case http.MethodPut:
		q.closed.Store(true)
		q.sse.Broadcast(quizEvent{Type: "reveal", Data: strconv.FormatInt(atomic.LoadInt64(&q.question), 10)})
	}
	w.WriteHeader(http.StatusNoContent)
}
