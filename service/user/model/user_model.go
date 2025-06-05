package model

type Users struct {
	Username *string
	Password *string
	Token    *string
	ID       uint `json:"id"`
}
