package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func main() {
	h := sseHandler{clients: &sync.Map{}}
	var question int64
	if start := os.Getenv("QUESTION"); start != "" {
		question, _ = strconv.ParseInt(start, 10, 64)
	}
	auth := authenticator{Secret: []byte(os.Getenv("SESSION_SECRET"))}
	q := &quizHandler{
		sse:      h,
		auth:     auth,
		question: question,
		closed:   new(atomic.Bool),
	}

	http.Handle("/events", h)
	http.Handle("/quiz", q)
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/admin/", q.Admin)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "index.html", nil)
	})
	if os.Getenv("NGROK_AUTHTOKEN") != "" {
		ctx := context.Background()
		l, err := ngrok.Listen(ctx,
			config.HTTPEndpoint(),
			ngrok.WithAuthtokenFromEnv(),
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("starting server on %s", l.URL())
		go http.Serve(l, nil)
	}
	// serve locally if no ngrok
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
