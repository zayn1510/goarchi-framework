package seeders

import "gorm.io/gorm"

type Seeder struct {
	Name string
	Up   func(*gorm.DB) error
	Down func(*gorm.DB) error
}

var AllSeeders = []Seeder{
}