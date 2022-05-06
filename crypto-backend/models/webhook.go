package models

type LimitMsg struct {
	UserId string `json:"userId"`
	Limit  Limit  `json:"limit"`
}
