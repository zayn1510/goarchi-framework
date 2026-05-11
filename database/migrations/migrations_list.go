package migrations

import "gorm.io/gorm"

type Migration struct {
	Name string
	Up   func(*gorm.DB) error
	Down func(*gorm.DB) error
}

var AllMigrations = []Migration{
}