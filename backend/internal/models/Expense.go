package models

import "time"

type Expense struct {
	Id      int       `json:"id"`
	Date    time.Time `json:"date"`
	GiftFor string    `json:"gift_for"`
	PupilId int       `json:"pupil_id"`
	Summ    int       `json:"summ"`
	Surname string    `json:"surname"`
}
