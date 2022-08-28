package models

import "gorm.io/gorm"

type Domain struct {
	gorm.Model
	Domain string
}
