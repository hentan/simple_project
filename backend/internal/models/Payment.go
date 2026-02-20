package models

type Payment struct {
	Id      int `json:"id" db:"id"`
	PupilId int `json:"pupil_id" db:"pupil_id"`
	Summ    int `json:"summ" db:"summ"`
}
