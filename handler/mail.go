package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	"github.com/gugabfigueiredo/dream-mail-go/service"
	log "github.com/gugabfigueiredo/tiny-go-log"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Handler struct {
	Service      service.IService
	Logger       *log.Logger
	MailingQueue chan models.Mail
	RetryQueue   chan models.Mail
}

func NewHandler(service service.IService, logger *log.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) HandleSend(w http.ResponseWriter, r *http.Request) {

	logger := h.Logger.C()

	mail, err := readMailFromRequest(r)
	if err != nil {
		logger.E("invalid or corrupted e-mail data", "err", err)
		http.Error(w, "invalid or corrupted e-mail data", http.StatusBadRequest)
		return
	}
	// queue mail for delivery
	h.Service.QueueMail(mail)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode([]byte(`{"status":  "OK", "message": "e-mail queued for delivery"}`)); err != nil {
		logger.E("error on json encoding", "err", err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}

	logger.I("e-mail queued for delivery")
}

// readMailFromRequest gets the email to be sent with all specs and unique ID, it returns error if mail is missing info
func readMailFromRequest(r *http.Request) (*models.Mail, error) {

	var mail models.Mail
	if err := json.NewDecoder(r.Body).Decode(&mail); err != nil {
		return &models.Mail{}, err
	}

	if ok, err := mail.Validate(); !ok {
		return &models.Mail{}, err
	}

	if mail.ID == "" {
		mail.ID = uuid.New().String()
	}

	return &mail, nil
}
