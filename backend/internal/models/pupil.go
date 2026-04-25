package models

type Pupil struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Surname     string `json:"surname" db:"surname"`
	ParentName  string `json:"parent_name" db:"parent_name"`
	ParentPhone string `json:"parent_phone" db:"parent_phone"`
}
