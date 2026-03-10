package models

import "time"

type Payment struct {
	Id      int       `json:"id" db:"id"`
	Date    time.Time `json:"date" db:"date"`
	PupilId int       `json:"pupil_id" db:"pupil_id"`
	Summ    int       `json:"summ" db:"summ"`
	Surname string    `json:"surname" db:"surname"`
	Purpose string    `json:"purpose" db:"purpose"`
}
