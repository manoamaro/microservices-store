package models

import "gorm.io/gorm"

type Flag struct {
	gorm.Model
	Name string
}
