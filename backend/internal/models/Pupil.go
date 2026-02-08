package models

type Pupil struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Surname     string `db:"surname"`
	ParentName  string `db:"parent_name"`
	ParentPhone string `db:"parent_phone"`
}
