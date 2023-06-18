package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gugabfigueiredo/dream-mail-go/env"
	"github.com/gugabfigueiredo/dream-mail-go/handler"
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/service"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"os"
	"time"
)

var Logger *log.Logger

func init() {
	envconfig.MustProcess("dmail", &env.Settings)

	Logger = log.New(env.Settings.Log)

	name, _ := os.Hostname()
	Logger = Logger.C("host", name)

}

func main() {

	// Start service
	mailService := service.NewService([]service.IProvider{
		service.NewSESProvider(env.Settings.Service.SES, Logger.C("provider", "ses")),
		service.NewSparkpostProvider(env.Settings.Service.Sparkpost, Logger.C("provider", "sparkpost")),
		service.NewSendgridProvider(env.Settings.Service.Sendgrid, Logger.C("provider", "sendgrid")),
		service.NewSMTPProvider(env.Settings.Service.SMTP, Logger.C("provider", "smtp")),
	}, Logger)

	// Handlers
	mailHandler := handler.NewHandler(mailService, Logger)

	// Start server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	r.Route(fmt.Sprintf("/%s", env.Settings.Server.Context), func(r chi.Router) {
		r.Post("/send", mailHandler.HandleSend)
	})

	http.Handle("/", r)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", env.Settings.Server.Port),
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	Logger.I("Starting server...", "port", env.Settings.Server.Port)

	if err := server.ListenAndServe(); err != nil {
		mailService.Quit()
		Logger.F("listen and serve died", "err", err)
	}
}
