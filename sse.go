package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type sseHandler struct {
	clients *sync.Map
}

type quizEvent struct {
	QuestionNumber int64
	Type           string
	From           string
	Data           string
}

func (s sseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusBadRequest)
		return
	}
	msgs := make(chan quizEvent)
	s.clients.Store(msgs, true)
	go func() {
		<-r.Context().Done()
		close(msgs)
		s.clients.Delete(msgs)
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for msg := range msgs {
		buf := new(bytes.Buffer)
		if err := templates.ExecuteTemplate(buf, msg.Type+".html", msg); err != nil {
			log.Printf("couldn't send message: %s", err)
		}
		data := strings.ReplaceAll(buf.String(), "\n", "\ndata: ")
		fmt.Fprintf(w, "data: %s\n\n", data)
		f.Flush()
	}
}

func (s sseHandler) Broadcast(event quizEvent) {
	s.clients.Range(func(client, _ any) bool {
		client.(chan quizEvent) <- event
		return true
	})
}
