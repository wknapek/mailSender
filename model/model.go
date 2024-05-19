package model

import (
	_ "gorm.io/gorm"
	"time"
)

type Email struct {
	Content    string    `json:"content"`
	Email      string    `json:"email"`
	InsertTime time.Time `json:"insert_time,omitempty" gorm:"autoCreateTime"`
	MailingID  int64     `json:"mailing_id" gorm:"index:primaryKey,unique"`
	Title      string    `json:"title"`
}

type Mailing struct {
	MailingID int64 `json:"mailing_id"`
}
