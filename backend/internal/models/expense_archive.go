package models

import "time"

type ExpenseArchive struct {
	ArchivedAt time.Time `json:"archived_at"`
	Operation  string    `json:"op"`
	ID         int       `json:"id"`
	Date       time.Time `json:"date"`
	Purpose    string    `json:"gift_for"`
	PupilID    int       `json:"pupil_id"`
	Amount     int       `json:"summ"`
	Surname    string    `json:"surname"`
}
