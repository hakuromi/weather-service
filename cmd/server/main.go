package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()     // роутер обрабатывет наши адреса
	r.Use(middleware.Logger) // мидлваре - для всех эндпоинтов
	r.Get("/get", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
