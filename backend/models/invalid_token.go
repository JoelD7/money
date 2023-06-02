package models

import "time"

type InvalidToken struct {
	Token       string    `json:"token,omitempty"`
	Expire      int64     `json:"expire,omitempty"`
	CreatedDate time.Time `json:"created_date,omitempty"`
}
