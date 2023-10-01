package models

type Version struct {
	Name    string `bun:"name,pk"`
	Version int32  `bun:"version,notnull"`
}
