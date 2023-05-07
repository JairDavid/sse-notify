package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"sse/sse"
)

func main() {
	r := chi.NewRouter()
	s := sse.New()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	port := os.Getenv("PORT")

	r.HandleFunc("/notification", s.Handler)
	r.Post("/message", func(writer http.ResponseWriter, request *http.Request) {
		var message string
		id := request.URL.Query().Get("id")

		if id == "" {
			http.Error(writer, "Missing ID", http.StatusInternalServerError)
			return
		}

		err := json.NewDecoder(request.Body).Decode(&message)
		if err != nil {
			fmt.Print(err)
		}
		s.BroadcastAll(message, id)
	})

	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}

}
