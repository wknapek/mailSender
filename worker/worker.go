package worker

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"vodeno/model"
)

type Worker struct {
	db *gorm.DB
}

func NewWorker(db *gorm.DB) *Worker {
	return &Worker{db: db}
}
func (work *Worker) Deleter() error {
	var emails []model.Email
	res := work.db.Find(&emails)
	if res.Error != nil {
		return res.Error
	}
	for _, email := range emails {
		res := work.db.Where("mailing_id = ?", email.MailingID).Delete(&model.Email{})
		if res.Error != nil {
			log.Error().Err(res.Error)
		}
	}
	return nil
}
