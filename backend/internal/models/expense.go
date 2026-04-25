package models

import "time"

type Expense struct {
	ID      int       `json:"id"`
	Date    time.Time `json:"date"`
	Purpose string    `json:"gift_for"`
	PupilID int       `json:"pupil_id"`
	Amount  int       `json:"summ"`
	Surname string    `json:"surname"`
}
