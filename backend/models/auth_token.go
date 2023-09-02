package models

import "time"

type AuthToken struct {
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiration"`
}
