package models

type Payment struct {
	Id      int    `json:"id"`
	Surname string `json:"surname"`
	Sum     string `json:"sum"`
}
