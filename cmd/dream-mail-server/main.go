package main

import (
	"github.com/gugabfigueiredo/dream-mail-go/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	mailHandler := handler.Handler{}

	r := chi.NewRouter()

	r.Route("/mail", func(r chi.Router) {
		r.Post("/send", mailHandler.HandleSend)
	})

	http.ListenAndServe(":3000", r)
}
