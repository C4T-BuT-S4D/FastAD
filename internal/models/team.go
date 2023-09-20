package models

import (
	"fmt"
)

type Team struct {
	ID      int
	Name    string
	Address string
	Token   string
	Labels  map[string]string
}

func (t *Team) String() string {
	return fmt.Sprintf("%s (%s)", t.Name, t.Address)
}
