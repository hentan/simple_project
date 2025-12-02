package models

import "time"

type Expense struct {
	Id      int       `json:"id"`
	Date    time.Time `json:"date"`
	GiftFor string    `json:"gift_for"`
	Surname string    `json:"surname"`
	Sum     string    `json:"sum"`
}
