package handler

import (
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/service"
	"net/http"
)

type Handler struct {
	Logger    *log.Logger
	Providers []service.IProvider
}

func (h *Handler) HandleSend(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Email Sent\n"))
}
