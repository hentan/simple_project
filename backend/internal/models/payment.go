package models

import "time"

type Payment struct {
	ID      int       `json:"id" db:"id"`
	Date    time.Time `json:"date" db:"date"`
	PupilID int       `json:"pupil_id" db:"pupil_id"`
	Amount  int       `json:"summ" db:"summ"`
	Surname string    `json:"surname" db:"surname"`
	Purpose string    `json:"purpose" db:"purpose"`
}
