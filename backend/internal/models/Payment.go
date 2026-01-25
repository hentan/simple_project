package models

type Payment struct {
	Id      int `json:"id"`
	PupilID int `json:"pupil_id"`
	Summ    int `json:"sum"`
}
