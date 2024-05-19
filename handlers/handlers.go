package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	gomail "gopkg.in/mail.v2"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"vodeno/model"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (han *Handler) CreateEmail(w http.ResponseWriter, r *http.Request) {
	email := model.Email{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg("error reading body")
		return
	}
	err = json.Unmarshal(reqBody, &email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg("error unmarshalling body")
		return
	}
	res := han.db.Create(&email)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		log.Error().Err(res.Error).Msg("error creating email")
		return
	}
	w.WriteHeader(http.StatusOK)
	replay := fmt.Sprintf("created email id:%d", email.MailingID)
	log.Info().Msg("mail created successfully")
	w.Write([]byte(replay))
}

func (han *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	emailId := model.Mailing{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg("error reading body")
		return
	}
	err = json.Unmarshal(reqBody, &emailId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg("error unmarshalling body")
		return
	}
	email := model.Email{}
	res := han.db.Model(&email).Where("mailing_id = ?", emailId.MailingID).Find(&email)
	if res.Error != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(res.Error).Msg("error finding email")
	}
	mailToSend := gomail.NewMessage()
	mailToSend.SetHeader("To", email.Email)
	mailToSend.SetHeader("Subject", email.Title)
	mailToSend.SetBody("text/plain", email.Content)
	mailToSend.SetHeader("From", "raven32@interia.pl")
	dialer := gomail.NewDialer("smtp.freesmtpservers.com", 25, "", "")
	err = dialer.DialAndSend(mailToSend)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msg("error sending email")
		return
	}
	han.deleteEmailByID(strconv.FormatInt(email.MailingID, 10))
	w.WriteHeader(http.StatusOK)
}

func (han *Handler) DeleteEmail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := han.deleteEmailByID(id)
	if err != nil {
		log.Error().Err(err).Msg("error deleting email")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (han *Handler) deleteEmailByID(id string) error {
	res := han.db.Where("mailing_id = ?", id).Delete(&model.Email{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("error deleting email")
		return res.Error
	}
	return nil
}
