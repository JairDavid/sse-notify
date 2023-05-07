package sse

import (
	"context"
	"fmt"
	"net/http"
)

const format = "event:%s\ndata:%s\n\n"

type Session struct {
	Id             string
	w              http.ResponseWriter
	flusher        http.Flusher
	requestContext context.Context
	Notify         chan string
}

func NewSession(id string, w http.ResponseWriter, flusher http.Flusher, requestContext context.Context) *Session {
	return &Session{
		Id:             id,
		w:              w,
		flusher:        flusher,
		requestContext: requestContext,
		Notify:         make(chan string),
	}
}

func (s *Session) ListenNotification() {
	for {
		select {
		case n := <-s.Notify:
			fmt.Print("got: ", n)
			fmt.Fprintf(s.w, format, "message", n)
			s.flusher.Flush()
		case <-s.requestContext.Done():
			//when the connection context is lost, it will be returned to handler and remove the session
			return
		}
	}
}
