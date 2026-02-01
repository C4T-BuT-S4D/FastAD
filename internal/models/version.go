package models

import "github.com/uptrace/bun"

type Version struct {
	bun.BaseModel `bun:"versions,alias:v"`

	Name    string `bun:"name,pk"`
	Version int    `bun:"version,notnull"`
}
