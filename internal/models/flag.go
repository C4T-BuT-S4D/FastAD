package models

type Flag struct {
	ID      int
	Flag    string
	Public  string
	Private string

	Round int

	TeamID    int
	ServiceID int
}
