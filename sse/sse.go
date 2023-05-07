package sse

import (
	"net/http"
	"sync"
)

type Server struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

func New() *Server {
	return &Server{
		sessions: map[string]*Session{},
	}
}

func (s *Server) BroadcastAll(message string, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, s := range s.sessions {
		if s.Id != id {
			s.Notify <- message
		}
	}
}

func (s *Server) RegisterSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.Id] = session

}

func (s *Server) RemoveSession(client string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, client)
}

func (s *Server) Handler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id := request.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	session := NewSession(id, w, flusher, request.Context())
	s.RegisterSession(session)
	session.ListenNotification()
	s.RemoveSession(session.Id)
}
