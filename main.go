package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main() {
	h := sseHandler{clients: &sync.Map{}}
	var question int64
	if start := os.Getenv("QUESTION"); start != "" {
		question, _ = strconv.ParseInt(start, 10, 64)
	}
	auth := authenticator{Secret: []byte(os.Getenv("SESSION_SECRET"))}
	q := &quizHandler{sse: h, auth: auth, question: question}

	http.Handle("/events", h)
	http.Handle("/quiz", q)
	http.HandleFunc("/admin/", auth.RestrictedHandler(q.Admin, "admin"))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
